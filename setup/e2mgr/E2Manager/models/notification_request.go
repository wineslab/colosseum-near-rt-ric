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
	"time"
	"unsafe"
)

type NotificationRequest struct {
	RanName       string
	Len           int
	Payload       []byte
	StartTime     time.Time
	TransactionId []byte
	msgSrc unsafe.Pointer
}

func NewNotificationRequest(ranName string, payload []byte, startTime time.Time, transactionId []byte, msgSrc unsafe.Pointer) *NotificationRequest {
	return &NotificationRequest{
		ranName,
		len(payload),
		payload,
		startTime,
		transactionId,
		msgSrc,
	}
}

func (r NotificationRequest) GetMsgSrc() unsafe.Pointer{
	return r.msgSrc
}
