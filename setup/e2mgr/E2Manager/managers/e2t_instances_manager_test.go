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

package managers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const E2TAddress = "10.10.2.15:9800"
const E2TAddress2 = "10.10.2.16:9800"
const PodName = "som_ pod_name"

func initE2TInstancesManagerTest(t *testing.T) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *E2TInstancesManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)
	return readerMock, writerMock, e2tInstancesManager
}

func TestAddNewE2TInstanceSaveE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "GetE2TAddresses")
}

func TestAddNewE2TInstanceGetE2TAddressesInternalFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "SaveE2TAddresses")
}

func TestAddNewE2TInstanceSaveE2TAddressesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	E2TAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(E2TAddresses, nil)
	E2TAddresses = append(E2TAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", E2TAddresses).Return(common.NewResourceNotFoundError(""))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.NotNil(t, err)
}

func TestAddNewE2TInstanceNoE2TAddressesSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, common.NewResourceNotFoundError(""))
	e2tAddresses = append(e2tAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", e2tAddresses).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TAddresses", e2tAddresses)
}

func TestAddNewE2TInstanceEmptyE2TAddressesSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddresses = append(e2tAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", e2tAddresses).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TAddresses", e2tAddresses)
}

func TestAddNewE2TInstanceExistingE2TAddressesSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	E2TAddresses := []string{"10.0.1.15:3030"}
	rnibReaderMock.On("GetE2TAddresses").Return(E2TAddresses, nil)
	E2TAddresses = append(E2TAddresses, E2TAddress)
	rnibWriterMock.On("SaveE2TAddresses", E2TAddresses).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress, PodName)
	assert.Nil(t, err)
}

func TestGetE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	var e2tInstance *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("for test")))
	res, err := e2tInstancesManager.GetE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestGetE2TInstanceSuccess(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address, PodName)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)
	res, err := e2tInstancesManager.GetE2TInstance(address)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestAddRanToInstanceGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AddRansToInstance(E2TAddress, []string{"test1"})
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestAddRanToInstanceSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AddRansToInstance(E2TAddress, []string{"test1"})
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestAddRanToInstanceSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tInstance := entities.NewE2TInstance(E2TAddress, PodName)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)

	updateE2TInstance := *e2tInstance
	updateE2TInstance.AssociatedRanList = append(updateE2TInstance.AssociatedRanList, "test1")

	rnibWriterMock.On("SaveE2TInstance", &updateE2TInstance).Return(nil)

	err := e2tInstancesManager.AddRansToInstance(E2TAddress, []string{"test1"})
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveRanFromInstanceGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))
	err := e2tInstancesManager.RemoveRanFromInstance("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestRemoveRanFromInstanceSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.RemoveRanFromInstance("test1", E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveRanFromInstanceSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance.AssociatedRanList = []string{"test0", "test1"}
	updatedE2TInstance := *e2tInstance
	updatedE2TInstance.AssociatedRanList = []string{"test0"}
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	rnibWriterMock.On("SaveE2TInstance", &updatedE2TInstance).Return(nil)

	err := e2tInstancesManager.RemoveRanFromInstance("test1", E2TAddress)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesGetE2TAddressesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(fmt.Errorf("for test")))
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertNotCalled(t, "GetE2TInstances")
}

func TestSelectE2TInstancesEmptyE2TAddressList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertNotCalled(t, "GetE2TInstances")
}

func TestSelectE2TInstancesGetE2TInstancesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	addresses := []string{E2TAddress}
	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances", addresses).Return([]*entities.E2TInstance{}, common.NewInternalError(fmt.Errorf("for test")))
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesEmptyE2TInstancesList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	addresses := []string{E2TAddress}
	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances", addresses).Return([]*entities.E2TInstance{}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Empty(t, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesNoActiveE2TInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	addresses := []string{E2TAddress, E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.ToBeDeleted
	e2tInstance2.AssociatedRanList = []string{"test4", "test5", "test6", "test7"}

	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances", addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.NotNil(t, err)
	assert.Equal(t, "", address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSelectE2TInstancesSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	addresses := []string{E2TAddress, E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4", "test5", "test6", "test7"}

	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances", addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	address, err := e2tInstancesManager.SelectE2TInstance()
	assert.Nil(t, err)
	assert.Equal(t, E2TAddress, address)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestActivateE2TInstanceSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.State == entities.Active })).Return(nil)

	err := e2tInstancesManager.SetE2tInstanceState(E2TAddress, entities.ToBeDeleted, entities.Active)
	assert.Nil(t, err)
	assert.Equal(t, entities.Active, e2tInstance1.State)
	rnibWriterMock.AssertExpectations(t)
}

func TestActivateE2TInstance_RnibError(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, common.NewInternalError(errors.New("for test")))

	err := e2tInstancesManager.SetE2tInstanceState(E2TAddress, entities.ToBeDeleted, entities.Active)
	assert.NotNil(t, err)
	rnibWriterMock.AssertExpectations(t)
}

func TestActivateE2TInstance_NoInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, e2managererrors.NewResourceNotFoundError())

	err := e2tInstancesManager.SetE2tInstanceState(E2TAddress, entities.ToBeDeleted, entities.Active)

	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestResetKeepAliveTimestampGetInternalFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address, PodName)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, common.NewInternalError(errors.New("Error")))
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	err := e2tInstancesManager.ResetKeepAliveTimestamp(address)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestAResetKeepAliveTimestampSaveInternalFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address, PodName)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(errors.New("Error")))

	err := e2tInstancesManager.ResetKeepAliveTimestamp(address)
	assert.NotNil(t, err)
}

func TestResetKeepAliveTimestampSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address, PodName)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	err := e2tInstancesManager.ResetKeepAliveTimestamp(address)
	assert.Nil(t, err)
	rnibReaderMock.AssertCalled(t, "GetE2TInstance", address)
	rnibWriterMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
}

func TestResetKeepAliveTimestampToBeDeleted(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address, PodName)
	e2tInstance.State = entities.ToBeDeleted
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)

	err := e2tInstancesManager.ResetKeepAliveTimestamp(address)
	assert.Nil(t, err)
	rnibReaderMock.AssertCalled(t, "GetE2TInstance", address)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestResetKeepAliveTimestampsForAllE2TInstancesGetE2TInstancesFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(errors.New("Error")))
	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestResetKeepAliveTimestampsForAllE2TInstancesNoInstances(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, nil)
	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestResetKeepAliveTimestampsForAllE2TInstancesNoActiveInstances(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.ToBeDeleted
	rnibReaderMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestResetKeepAliveTimestampsForAllE2TInstancesOneActiveInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.ToBeDeleted
	rnibReaderMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()
	rnibWriterMock.AssertNumberOfCalls(t, "SaveE2TInstance",1)
}

func TestResetKeepAliveTimestampsForAllE2TInstancesSaveE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.ToBeDeleted
	rnibReaderMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(errors.New("Error")))
	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()
	rnibWriterMock.AssertNumberOfCalls(t, "SaveE2TInstance",1)
}

func TestRemoveE2TInstanceSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddressesNew := []string{E2TAddress2}
	rnibWriterMock.On("SaveE2TAddresses", e2tAddressesNew).Return(nil)

	err := e2tInstancesManager.RemoveE2TInstance(E2TAddress)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveE2TInstanceRnibErrorInRemoveInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("RemoveE2TInstance", E2TAddress).Return(e2managererrors.NewRnibDbError())

	err := e2tInstancesManager.RemoveE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveE2TInstanceRnibErrorInGetAddresses(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	var e2tAddresses []string
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, e2managererrors.NewRnibDbError())

	err := e2tInstancesManager.RemoveE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestRemoveE2TInstanceRnibErrorInSaveAddresses(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddressesNew := []string{E2TAddress2}
	rnibWriterMock.On("SaveE2TAddresses", e2tAddressesNew).Return(e2managererrors.NewRnibDbError())

	err := e2tInstancesManager.RemoveE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestSetE2tInstanceStateCurrentStateHasChanged(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance.State = entities.Active

	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)

	err := e2tInstancesManager.SetE2tInstanceState(E2TAddress, entities.ToBeDeleted, entities.Active)
	assert.NotNil(t, err)
	assert.IsType(t, e2managererrors.NewInternalError(), err)
	rnibReaderMock.AssertExpectations(t)
}

func TestSetE2tInstanceStateErrorInSaveE2TInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	e2tInstance := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance.State = entities.ToBeDeleted
	rnibReaderMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for testing")))

	err := e2tInstancesManager.SetE2tInstanceState(E2TAddress, entities.ToBeDeleted, entities.Active)
	assert.NotNil(t, err)
	assert.IsType(t, &common.InternalError{}, err)
	rnibReaderMock.AssertExpectations(t)
}

func TestClearRansOfAllE2TInstancesEmptyList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	E2TAddresses := []string{}
	rnibReaderMock.On("GetE2TAddresses").Return(E2TAddresses, nil)
	err := e2tInstancesManager.ClearRansOfAllE2TInstances()
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestClearRansOfAllE2TInstancesErrorInSaveE2TInstance(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	addresses := []string{E2TAddress, E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4", "test5", "test6", "test7"}

	rnibReaderMock.On("GetE2TAddresses").Return(addresses, nil)
	rnibReaderMock.On("GetE2TInstances", addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress})).Return(common.NewInternalError(fmt.Errorf("for testing")))
	rnibWriterMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress2})).Return(nil)
	err := e2tInstancesManager.ClearRansOfAllE2TInstances()
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}
