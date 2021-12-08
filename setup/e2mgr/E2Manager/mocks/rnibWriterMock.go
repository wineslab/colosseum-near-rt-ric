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

type RnibWriterMock struct {
	mock.Mock
}

func (rnibWriterMock *RnibWriterMock) SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error {
	args := rnibWriterMock.Called(nbIdentity, nb)

	errArg := args.Get(0)

	if errArg != nil {
		return errArg.(error)
	}

	return nil
}

func (rnibWriterMock *RnibWriterMock) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {
	args := rnibWriterMock.Called(nodebInfo)

	errArg := args.Get(0)

	if errArg != nil {
		return errArg.(error)
	}

	return nil
}

func (rnibWriterMock *RnibWriterMock) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {
	args := rnibWriterMock.Called(inventoryName, ranLoadInformation)

	errArg := args.Get(0)

	if errArg != nil {
		return errArg.(error)
	}

	return nil
}

func (rnibWriterMock *RnibWriterMock) SaveE2TInstance(e2tInstance *entities.E2TInstance) error {
	args := rnibWriterMock.Called(e2tInstance)

	return args.Error(0)
}

func (rnibWriterMock *RnibWriterMock) SaveE2TAddresses(addresses []string) error {
	args := rnibWriterMock.Called(addresses)

	return args.Error(0)
}

func (rnibWriterMock *RnibWriterMock) RemoveE2TInstance(address string) error {
	args := rnibWriterMock.Called(address)

	return args.Error(0)
}

func (rnibWriterMock *RnibWriterMock) UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error {
	args := rnibWriterMock.Called(nodebInfo, servedNrCells)
	return args.Error(0)
}

func (rnibWriterMock *RnibWriterMock) RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error {
	args := rnibWriterMock.Called(inventoryName, servedNrCells)
	return args.Error(0)
}

