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


package enums

import (
	"encoding/json"
	"strconv"
)

type MessageDirection int32

var messageDirectionEnumName = map[int32]string{
	0: "UNKNOWN_MESSAGE_DIRECTION",
	1: "RAN_TO_RIC",
	2: "RIC_TO_RAN",
}

const (
	UNKNOWN_MESSAGE_DIRECTION MessageDirection = 0
	RAN_TO_RIC                MessageDirection = 1
	RIC_TO_RAN                MessageDirection = 2
)

func (md MessageDirection) String() string {
	s, ok := messageDirectionEnumName[int32(md)]
	if ok {
		return s
	}
	return strconv.Itoa(int(md))
}

func (md MessageDirection) MarshalJSON() ([]byte, error) {
	_, ok := messageDirectionEnumName[int32(md)]

	if !ok {
		return nil,&json.UnsupportedValueError{}
	}

	v:= int32(md)
	return json.Marshal(v)
}
