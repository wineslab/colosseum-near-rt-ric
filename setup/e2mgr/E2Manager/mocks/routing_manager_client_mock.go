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

package mocks

import (
	"github.com/stretchr/testify/mock"
)

type RoutingManagerClientMock struct {
	mock.Mock
}

func (m *RoutingManagerClientMock) AddE2TInstance(e2tAddress string) error {

	args := m.Called(e2tAddress)
	return args.Error(0)
}

func (m *RoutingManagerClientMock) AssociateRanToE2TInstance(e2tAddress string, ranName string) error {

	args := m.Called(e2tAddress, ranName)
	return args.Error(0)
}

func (m *RoutingManagerClientMock) DissociateRanE2TInstance(e2tAddress string, ranName string) error {

	args := m.Called(e2tAddress, ranName)
	return args.Error(0)
}

func (m *RoutingManagerClientMock) DissociateAllRans(e2tAddresses []string) error {

args := m.Called(e2tAddresses)
return args.Error(0)
}

func (m *RoutingManagerClientMock) DeleteE2TInstance(e2tAddress string, ransToBeDissociated []string) error {

	args := m.Called(e2tAddress, ransToBeDissociated)
	return args.Error(0)
}