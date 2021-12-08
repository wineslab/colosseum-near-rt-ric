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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/mock"
)

type RnibReaderMock struct {
	mock.Mock
}

func (m *RnibReaderMock) GetNodeb(inventoryName string) (*entities.NodebInfo, error) {
	args := m.Called(inventoryName)

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).(*entities.NodebInfo), errArg.(error);
	}

	return args.Get(0).(*entities.NodebInfo), nil
}

func (m *RnibReaderMock) GetNodebByGlobalNbId(nodeType entities.Node_Type, globalNbId *entities.GlobalNbId) (*entities.NodebInfo, error) {
	args := m.Called(nodeType, globalNbId)

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).(*entities.NodebInfo), errArg.(error);
	}

	return args.Get(0).(*entities.NodebInfo), nil
}

func (m *RnibReaderMock) GetCellList(inventoryName string) (*entities.Cells, error) {
	args := m.Called(inventoryName)

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).(*entities.Cells), errArg.(error);
	}

	return args.Get(0).(*entities.Cells), nil
}

func (m *RnibReaderMock) GetListGnbIds() ([]*entities.NbIdentity, error) {
	args := m.Called()

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).([]*entities.NbIdentity), errArg.(error);
	}

	return args.Get(0).([]*entities.NbIdentity), nil
}

func (m *RnibReaderMock) GetListEnbIds() ([]*entities.NbIdentity, error) {
	args := m.Called()

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).([]*entities.NbIdentity), errArg.(error);
	}

	return args.Get(0).([]*entities.NbIdentity), nil

}

func (m *RnibReaderMock) GetCountGnbList() (int, error) {
	args := m.Called()

	errArg := args.Get(1);
	if errArg != nil {
		return args.Int(0), errArg.(error);
	}

	return args.Int(0), nil

}

func (m *RnibReaderMock) GetCell(inventoryName string, pci uint32) (*entities.Cell, error) {
	args := m.Called(inventoryName, pci)

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).(*entities.Cell), errArg.(error);
	}

	return args.Get(0).(*entities.Cell), nil
}

func (m *RnibReaderMock) GetCellById(cellType entities.Cell_Type, cellId string) (*entities.Cell, error) {
	args := m.Called(cellType, cellId)

	errArg := args.Get(1);
	if errArg != nil {
		return args.Get(0).(*entities.Cell), errArg.(error);
	}

	return args.Get(0).(*entities.Cell), nil
}

func (m *RnibReaderMock) GetListNodebIds() ([]*entities.NbIdentity, error) {
	args := m.Called()

	errArg := args.Get(1)

	if errArg != nil {
		return args.Get(0).([]*entities.NbIdentity), errArg.(error)
	}

	return args.Get(0).([]*entities.NbIdentity), nil
}

func (m *RnibReaderMock) GetRanLoadInformation(inventoryName string) (*entities.RanLoadInformation, error) {
	args := m.Called()

	errArg := args.Get(1)

	if errArg != nil {
		return args.Get(0).(*entities.RanLoadInformation), errArg.(error)
	}

	return args.Get(0).(*entities.RanLoadInformation), nil
}

func (m *RnibReaderMock) GetE2TInstance(e2taddress string) (*entities.E2TInstance, error) {
	args := m.Called(e2taddress)
	return args.Get(0).(*entities.E2TInstance), args.Error(1)
}

func (m *RnibReaderMock) GetE2TInstances(addresses []string) ([]*entities.E2TInstance, error) {
	args := m.Called(addresses)
	return args.Get(0).([]*entities.E2TInstance), args.Error(1)
}

func (m *RnibReaderMock) GetE2TAddresses() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
