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

package rmr

import (
	"strconv"
)

// RmrService holds an instance of RMR messenger as well as its configuration
type Service struct {
	messenger *Messenger
}

// NewRmrService instantiates a new Rmr service instance
func NewService(rmrConfig Config, messenger Messenger) *Service {
	return &Service{
		messenger: messenger.Init("tcp:"+strconv.Itoa(rmrConfig.Port), rmrConfig.MaxMsgSize, rmrConfig.MaxRetries, rmrConfig.Flags),
	}
}

func (r *Service) SendMessage(messageType int, ranName string, msg []byte, transactionId []byte) (*MBuf, error) {
	mbuf := NewMBuf(messageType, len(msg), msg, transactionId)
	mbuf.Meid = ranName
	return (*r.messenger).SendMsg(mbuf)
}

func (r *Service) RecvMessage() (*MBuf, error) {
	return (*r.messenger).RecvMsg()
}

func (r *Service) CloseContext() {
	(*r.messenger).Close()

}
