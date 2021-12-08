//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package dispatcher

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
	"xappmock/enums"
	"xappmock/logger"
	"xappmock/models"
	"xappmock/rmr"
	"xappmock/sender"
)

// Id -> Command
var configuration = make(map[string]*models.JsonCommand)

// Rmr Message Id -> Command
var waitForRmrMessageType = make(map[int]*models.JsonCommand)

func addRmrMessageToWaitFor(rmrMessageToWaitFor string, command models.JsonCommand) error {
	rmrMsgId, err := rmr.MessageIdToUint(rmrMessageToWaitFor)

	if err != nil {
		return errors.New(fmt.Sprintf("invalid rmr message id: %s", rmrMessageToWaitFor))
	}

	waitForRmrMessageType[int(rmrMsgId)] = &command
	return nil
}

type Dispatcher struct {
	rmrService    *rmr.Service
	processResult models.ProcessResult
	logger        *logger.Logger
	jsonSender    *sender.JsonSender
}

func (d *Dispatcher) GetProcessResult() models.ProcessResult {
	return d.processResult
}

func New(logger *logger.Logger, rmrService *rmr.Service, jsonSender *sender.JsonSender) *Dispatcher {
	return &Dispatcher{
		rmrService: rmrService,
		logger:     logger,
		jsonSender: jsonSender,
	}
}

func (d *Dispatcher) JsonCommandsDecoderCB(cmd models.JsonCommand) error {
	if len(cmd.Id) == 0 {
		return errors.New(fmt.Sprintf("invalid cmd, no id"))
	}
	configuration[cmd.Id] = &cmd
	return nil

	//	if len(cmd.ReceiveCommandId) == 0 {
	//		return nil
	//	}
	//
	//	return addRmrMessageToWaitFor(cmd.ReceiveCommandId, cmd)
}

func (d *Dispatcher) sendNoRepeat(command models.JsonCommand) error {

	if enums.CommandAction(command.Action) == enums.SendRmrMessage && d.processResult.StartTime == nil {
		now := time.Now()
		d.processResult.StartTime = &now
	}

	err := d.jsonSender.SendJsonRmrMessage(command, nil, d.rmrService)

	if err != nil {
		d.logger.Errorf("#Dispatcher.sendNoRepeat - error sending rmr message: %s", err)
		d.processResult.Err = err
		d.processResult.Stats.SentErrorCount.Inc()
		return err
	}

	d.processResult.Stats.SentCount.Inc()
	return nil
}

func (d *Dispatcher) sendWithRepeat(ctx context.Context, command models.JsonCommand) {

	if enums.CommandAction(command.Action) == enums.SendRmrMessage && d.processResult.StartTime == nil {
		now := time.Now()
		d.processResult.StartTime = &now
	}

	for repeatCount := command.RepeatCount; repeatCount > 0; repeatCount-- {

		select {
		case <-ctx.Done():
			return
		default:
		}

		err := d.jsonSender.SendJsonRmrMessage(command, nil, d.rmrService)

		if err != nil {
			d.logger.Errorf("#Dispatcher.sendWithRepeat - error sending rmr message: %s", err)
			d.processResult.Stats.SentErrorCount.Inc()
			continue
		}

		d.processResult.Stats.SentCount.Inc()
		time.Sleep(time.Duration(command.RepeatDelayInMs) * time.Millisecond)
	}
}

func getReceiveRmrMessageType(receiveCommandId string) (string, error) {
	command, ok := configuration[receiveCommandId]

	if !ok {
		return "", errors.New(fmt.Sprintf("invalid receive command id: %s", receiveCommandId))
	}

	if len(command.RmrMessageType) == 0 {
		return "", errors.New(fmt.Sprintf("missing RmrMessageType for command id: %s", receiveCommandId))
	}

	return command.RmrMessageType, nil
}

func (d *Dispatcher) sendHandler(ctx context.Context, sendAndReceiveWg *sync.WaitGroup, command models.JsonCommand) {

	defer sendAndReceiveWg.Done()
	var listenAndHandleWg sync.WaitGroup

	if len(command.ReceiveCommandId) > 0 {
		rmrMessageToWaitFor, err := getReceiveRmrMessageType(command.ReceiveCommandId)

		if err != nil {
			d.processResult.Err = err
			return
		}

		err = addRmrMessageToWaitFor(rmrMessageToWaitFor, command)

		if err != nil {
			d.processResult.Err = err
			return
		}

		listenAndHandleWg.Add(1)
		go d.listenAndHandle(ctx, &listenAndHandleWg, command)
	}

	if command.RepeatCount == 0 {
		err := d.sendNoRepeat(command)

		if err != nil {
			return
		}

	} else {
		d.sendWithRepeat(ctx, command)
	}

	if len(command.ReceiveCommandId) > 0 {
		listenAndHandleWg.Wait()
	}
}

func (d *Dispatcher) receiveHandler(ctx context.Context, sendAndReceiveWg *sync.WaitGroup, command models.JsonCommand) {

	defer sendAndReceiveWg.Done()

	err := addRmrMessageToWaitFor(command.RmrMessageType, command)

	if err != nil {
		d.processResult.Err = err
		return
	}

	var listenAndHandleWg sync.WaitGroup
	listenAndHandleWg.Add(1) // this is due to the usage of listenAndHandle as a goroutine in the sender case
	d.listenAndHandle(ctx, &listenAndHandleWg, command)
}

func getMergedCommand(cmd *models.JsonCommand) (models.JsonCommand, error) {
	var command models.JsonCommand
	if len(cmd.Id) == 0 {
		return command, errors.New(fmt.Sprintf("invalid command, no id"))
	}

	command = *cmd

	conf, ok := configuration[cmd.Id]

	if ok {
		command = *conf
		mergeConfigurationAndCommand(&command, cmd)
	}

	return command, nil
}

func (d *Dispatcher) ProcessJsonCommand(ctx context.Context, cmd *models.JsonCommand) {

	command, err := getMergedCommand(cmd)

	if err != nil {
		d.processResult.Err = err
		return
	}

	var sendAndReceiveWg sync.WaitGroup

	commandAction := enums.CommandAction(command.Action)

	switch commandAction {

	case enums.SendRmrMessage:
		sendAndReceiveWg.Add(1)
		go d.sendHandler(ctx, &sendAndReceiveWg, command)
	case enums.ReceiveRmrMessage:
		sendAndReceiveWg.Add(1)
		go d.receiveHandler(ctx, &sendAndReceiveWg, command)
	default:
		d.processResult = models.ProcessResult{Err: errors.New(fmt.Sprintf("invalid command action %s", command.Action))}
		return
	}

	sendAndReceiveWg.Wait()
}

func getResponseCommand(command models.JsonCommand) (*models.JsonCommand, error) {
	responseCommand, ok := configuration[command.SendCommandId]

	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid SendCommandId %s", command.SendCommandId))
	}

	return responseCommand, nil
}

func (d *Dispatcher) listenAndHandleNoRepeat(ctx context.Context, command models.JsonCommand) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mbuf, err := d.rmrService.RecvMessage()

		if err != nil {
			d.logger.Errorf("#Dispatcher.listenAndHandleNoRepeat - error receiving message: %s", err)
			d.processResult.Err = err
			d.processResult.Stats.ReceivedErrorCount.Inc()
			return
		}

		if enums.CommandAction(command.Action) == enums.ReceiveRmrMessage && d.processResult.StartTime == nil {
			now := time.Now()
			d.processResult.StartTime = &now
		}

		messageInfo := models.NewMessageInfo(mbuf.MType, mbuf.Meid, mbuf.Payload, mbuf.XAction)

		_, ok := waitForRmrMessageType[mbuf.MType]

		if !ok {
			d.logger.Infof("#Dispatcher.listenAndHandleNoRepeat - received unexpected msg: %s", messageInfo)
			d.processResult.Stats.ReceivedUnexpectedCount.Inc()
			continue
		}

		d.logger.Infof("#Dispatcher.listenAndHandleNoRepeat - received expected msg: %s", messageInfo)
		d.processResult.Stats.ReceivedExpectedCount.Inc()

		if len(command.SendCommandId) > 0 {
			responseCommand, err := getResponseCommand(command)

			if err != nil {
				d.processResult.Err = err
				return
			}

			_ = d.sendNoRepeat(*responseCommand)
		}

		return
	}
}

func (d *Dispatcher) listenAndHandleWithRepeat(ctx context.Context, command models.JsonCommand) {

	var responseCommand *models.JsonCommand

	if len(command.SendCommandId) > 0 {
		var err error
		responseCommand, err = getResponseCommand(command)

		if err != nil {
			d.processResult.Err = err
			return
		}
	}

	for d.processResult.Stats.ReceivedExpectedCount.Load() < int32(command.RepeatCount) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mbuf, err := d.rmrService.RecvMessage()

		if err != nil {
			d.logger.Errorf("#Dispatcher.listenAndHandleWithRepeat - error receiving message: %s", err)
			d.processResult.Stats.ReceivedErrorCount.Inc()
			continue
		}

		if enums.CommandAction(command.Action) == enums.ReceiveRmrMessage && d.processResult.StartTime == nil {
			now := time.Now()
			d.processResult.StartTime = &now
		}

		messageInfo := models.NewMessageInfo(mbuf.MType, mbuf.Meid, mbuf.Payload, mbuf.XAction)

		_, ok := waitForRmrMessageType[mbuf.MType]

		if !ok {
			d.logger.Infof("#Dispatcher.listenAndHandleWithRepeat - received unexpected msg: %s", messageInfo)
			d.processResult.Stats.ReceivedUnexpectedCount.Inc()
			continue
		}

		d.logger.Infof("#Dispatcher.listenAndHandleWithRepeat - received expected msg: %s", messageInfo)
		d.processResult.Stats.ReceivedExpectedCount.Inc()

		if responseCommand != nil {
			_ = d.sendNoRepeat(*responseCommand) // TODO: goroutine? + error handling
		}
	}
}

func (d *Dispatcher) listenAndHandle(ctx context.Context, listenAndHandleWg *sync.WaitGroup, command models.JsonCommand) {

	defer listenAndHandleWg.Done()

	if command.RepeatCount == 0 {
		d.listenAndHandleNoRepeat(ctx, command)
		return
	}

	d.listenAndHandleWithRepeat(ctx, command)
}

func mergeConfigurationAndCommand(conf *models.JsonCommand, cmd *models.JsonCommand) {
	nFields := reflect.Indirect(reflect.ValueOf(cmd)).NumField()

	for i := 0; i < nFields; i++ {
		if fieldValue := reflect.Indirect(reflect.ValueOf(cmd)).Field(i); fieldValue.IsValid() {
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.Len() > 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if fieldValue.Int() != 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Bool:
				if fieldValue.Bool() {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			case reflect.Float64, reflect.Float32:
				if fieldValue.Float() != 0 {
					reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
				}
			default:
				reflect.Indirect(reflect.ValueOf(conf)).Field(i).Set(fieldValue)
			}
		}
	}
}
