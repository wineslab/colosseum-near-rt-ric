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

package mocks

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type NodebControllerMock struct {
	mock.Mock
}

func (c *NodebControllerMock) GetNodeb(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	ranName := vars["ranName"]

	writer.Write([]byte(ranName))
	c.Called()
}

func (c *NodebControllerMock) GetNodebIdList(writer http.ResponseWriter, r *http.Request) {
	c.Called()
}

func (c *NodebControllerMock) Shutdown(writer http.ResponseWriter, r *http.Request) {
	c.Called()
}

func (c *NodebControllerMock) X2Reset(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	ranName := vars["ranName"]

	writer.Write([]byte(ranName))

	c.Called()
}

func (c *NodebControllerMock) X2Setup(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	c.Called()
}

func (c *NodebControllerMock) EndcSetup(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	c.Called()
}

func (c *NodebControllerMock) UpdateGnb(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	c.Called()
}
