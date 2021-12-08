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

package models

import (
	"fmt"
	"time"
)

// TODO: message command id / source / dest

type MessageInfo struct {
	MessageTimestamp int64  `json:"messageTimestamp"`
	MessageType      int    `json:"messageType"`
	Meid             string `json:"meid"`
	Payload          []byte `json:"payload"`
	TransactionId    string `json:"transactionId"`
}

func NewMessageInfo(messageType int, meid string, payload []byte, transactionId []byte) MessageInfo {
	return MessageInfo{
		MessageTimestamp: time.Now().Unix(),
		MessageType:      messageType,
		Meid:             meid,
		Payload:          payload,
		TransactionId:    string(transactionId),
	}
}

func (mi MessageInfo) String() string {
	return fmt.Sprintf("message timestamp: %d | message type: %d | meid: %s | payload: %x | transaction id: %s",
		mi.MessageTimestamp, mi.MessageType, mi.Meid, mi.Payload, mi.TransactionId)
}
