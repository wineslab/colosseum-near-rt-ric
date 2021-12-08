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

type RoutingManagerE2TData struct {
	E2TAddress  string   `json:"E2TAddress"`
	RanNamelist []string `json:"ranNamelist,omitempty"`
}

func NewRoutingManagerE2TData (e2tAddress string, ranNameList ...string) *RoutingManagerE2TData {
	data := &RoutingManagerE2TData{
		E2TAddress: e2tAddress,
	}

	if len(ranNameList) == 0 {
		return data
	}

	for _, ranName := range ranNameList {
		data.RanNamelist = append(data.RanNamelist, ranName)
	}

	return data
}
