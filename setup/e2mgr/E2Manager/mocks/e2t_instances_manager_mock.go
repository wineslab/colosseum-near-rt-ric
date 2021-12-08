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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/mock"
)

type E2TInstancesManagerMock struct {
	mock.Mock
}

func (m *E2TInstancesManagerMock) GetE2TInstance(e2tAddress string) (*entities.E2TInstance, error) {
	args := m.Called(e2tAddress)

	return args.Get(0).(*entities.E2TInstance), args.Error(1)
}

func (m *E2TInstancesManagerMock) AddE2TInstance(e2tInstanceAddress string, podName string) error {
	args := m.Called(e2tInstanceAddress, podName)
	return args.Error(0)
}

func (m *E2TInstancesManagerMock) RemoveE2TInstance(e2tAddress string) error {
	args := m.Called(e2tAddress)
	return args.Error(0)
}

func (m *E2TInstancesManagerMock) SelectE2TInstance() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *E2TInstancesManagerMock) AddRansToInstance(e2tAddress string, ranNames []string) error {
	args := m.Called(e2tAddress, ranNames)
	return args.Error(0)

}

func (m *E2TInstancesManagerMock) RemoveRanFromInstance(ranName string, e2tAddress string) error {
	args := m.Called(ranName, e2tAddress)
	return args.Error(0)

}

func (m *E2TInstancesManagerMock) GetE2TInstances() ([]*entities.E2TInstance, error) {
	args := m.Called()

	return args.Get(0).([]*entities.E2TInstance), args.Error(1)
}

func (m *E2TInstancesManagerMock) GetE2TInstancesNoLogs() ([]*entities.E2TInstance, error) {
	args := m.Called()

	return args.Get(0).([]*entities.E2TInstance), args.Error(1)
}

func (m *E2TInstancesManagerMock) ResetKeepAliveTimestamp(e2tAddress string) error {
	args := m.Called(e2tAddress)
	return args.Error(0)

}

func (m *E2TInstancesManagerMock) SetE2tInstanceState(e2tAddress string, currentState entities.E2TInstanceState, newState entities.E2TInstanceState) error {
	args := m.Called(e2tAddress, currentState, newState)
	return args.Error(0)
}

func (m *E2TInstancesManagerMock) GetE2TAddresses() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *E2TInstancesManagerMock) ClearRansOfAllE2TInstances() error {
	args := m.Called()
	return args.Error(0)
}
