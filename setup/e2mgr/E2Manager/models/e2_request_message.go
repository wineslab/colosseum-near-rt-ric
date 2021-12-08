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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package models

import (
	"e2mgr/logger"
	"fmt"
)

type E2RequestMessage struct {
	transactionId string
	ranIp         string
	ranPort       uint16
	ranName       string
	payload       []byte
}

func (e2RequestMessage E2RequestMessage) RanName() string {
	return e2RequestMessage.ranName
}

func (e2RequestMessage E2RequestMessage) TransactionId() string {
	return e2RequestMessage.transactionId
}

func NewE2RequestMessage(transactionId string, ranIp string, ranPort uint16, ranName string, payload []byte) *E2RequestMessage {
	return &E2RequestMessage{transactionId: transactionId, ranIp: ranIp, ranPort: ranPort, ranName: ranName, payload: payload}
}

// TODO: this shouldn't receive logger
func (e2RequestMessage E2RequestMessage) GetMessageAsBytes(logger *logger.Logger) []byte {
	messageStringWithoutPayload := fmt.Sprintf("%s|%d|%s|%d|", e2RequestMessage.ranIp, e2RequestMessage.ranPort, e2RequestMessage.ranName, len(e2RequestMessage.payload))
	logger.Debugf("#e2_request_message.GetMessageAsBytes - messageStringWithoutPayload: %s", messageStringWithoutPayload)
	messageBytesWithoutPayload := []byte(messageStringWithoutPayload)
	return append(messageBytesWithoutPayload, e2RequestMessage.payload...)
}
