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
	"e2mgr/e2managererrors"
	"encoding/json"
)

type E2TInstancesResponse []*E2TInstanceResponseModel

type E2TInstanceResponseModel struct {
	E2TAddress string   `json:"e2tAddress"`
	RanNames   []string `json:"ranNames"`
}

func NewE2TInstanceResponseModel(e2tAddress string, ranNames []string) *E2TInstanceResponseModel {
	return &E2TInstanceResponseModel{
		E2TAddress: e2tAddress,
		RanNames:   ranNames,
	}
}

func (response E2TInstancesResponse) Marshal() ([]byte, error) {

	data, err := json.Marshal(response)

	if err != nil {
		return nil, e2managererrors.NewInternalError()
	}

	return data, nil

}
