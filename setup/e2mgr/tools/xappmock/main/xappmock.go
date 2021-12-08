//
// Copyright 2019 AT&T Intellectual Property
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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
	"xappmock/dispatcher"
	"xappmock/frontend"
	"xappmock/logger"
	"xappmock/rmr"
	"xappmock/sender"
)

const (
	ENV_RMR_PORT     = "RMR_PORT"
	RMR_PORT_DEFAULT = 5001
)

var rmrService *rmr.Service

func main() {

	logLevel, _ := logger.LogLevelTokenToLevel("info")
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}

	var rmrContext *rmr.Context
	var rmrConfig = rmr.Config{Port: RMR_PORT_DEFAULT, MaxMsgSize: rmr.RMR_MAX_MSG_SIZE, MaxRetries: 10, Flags: 0}

	if port, err := strconv.ParseUint(os.Getenv(ENV_RMR_PORT), 10, 16); err == nil {
		rmrConfig.Port = int(port)
	} else {
		logger.Infof("#main - %s: %s, using default (%d).", ENV_RMR_PORT, err, RMR_PORT_DEFAULT)
	}

	rmrService = rmr.NewService(rmrConfig, rmrContext)
	jsonSender := sender.NewJsonSender(logger)
	dispatcherDesc := dispatcher.New(logger, rmrService, jsonSender)

	/* Load configuration file*/
	err = frontend.ProcessConfigurationFile("resources", "conf", ".json",
		func(data []byte) error {
			return frontend.JsonCommandsDecoder(data, dispatcherDesc.JsonCommandsDecoderCB)
		})

	if err != nil {
		logger.Errorf("#main - processing error: %s", err)
		os.Exit(1)
	}

	logger.Infof("#main - xApp Mock is up and running...")

	flag.Parse()
	cmd := flag.Arg(0) /*first remaining argument after flags have been processed*/

	command, err := frontend.DecodeJsonCommand([]byte(cmd))

	if err != nil {
		logger.Errorf("#main - command decoding error: %s", err)
		rmrService.CloseContext()
		logger.Infof("#main - xApp Mock is down")
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		logger.Infof("system call:%+v", oscall)
		cancel()
		rmrService.CloseContext()
	}()

	dispatcherDesc.ProcessJsonCommand(ctx, command)
	pr := dispatcherDesc.GetProcessResult()

	if pr.Err != nil {
		logger.Errorf("#main - command processing Error: %s", pr.Err)
	}

	if pr.StartTime != nil {
		processElapsedTimeInMs := float64(time.Since(*pr.StartTime)) / float64(time.Millisecond)
		logger.Infof("#main - processing (sending/receiving) messages took %.2f ms", processElapsedTimeInMs)

	}
	logger.Infof("#main - process stats: %s", pr.Stats)

	rmrService.CloseContext() // TODO: called twice
	logger.Infof("#main - xApp Mock is down")
}
