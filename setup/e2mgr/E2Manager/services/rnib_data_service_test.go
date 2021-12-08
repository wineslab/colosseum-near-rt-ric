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


package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

func setupRnibDataServiceTest(t *testing.T) (*rNibDataService, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	return setupRnibDataServiceTestWithMaxAttempts(t, 3)
}

func setupRnibDataServiceTestWithMaxAttempts(t *testing.T, maxAttempts int) (*rNibDataService, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: maxAttempts}

	readerMock := &mocks.RnibReaderMock{}


	writerMock := &mocks.RnibWriterMock{}


	rnibDataService := NewRnibDataService(logger, config, readerMock, writerMock)
	assert.NotNil(t, rnibDataService)

	return rnibDataService, readerMock, writerMock
}

func TestSuccessfulSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 3)
}

func TestNonConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection failure")}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestSuccessfulUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(nil)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnFailureUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(mockErr)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 3)
}

func TestSuccessfulSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(nil)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 1)
}

func TestConnFailureSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(mockErr)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 3)
}

func TestSuccessfulGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	nodebInfo := &entities.NodebInfo{}
	readerMock.On("GetNodeb", invName).Return(nodebInfo, nil)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
	assert.Equal(t, nodebInfo, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeb, res)
}

func TestSuccessfulGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	nodeIds := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.Equal(t, nodeIds, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeIds, res)
}

func TestConnFailureTwiceGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeIds, res)

	res2, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeb, res2)
}

func TestConnFailureWithAnotherConfig(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTestWithMaxAttempts(t, 5)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 5)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeIds, res)
}

func TestPingRnibConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.False(t, res)
}

func TestPingRnibOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}

func TestPingRnibOkOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}

//func TestConnFailureThenSuccessGetNodebIdList(t *testing.T) {
//	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)
//
//	var nilNodeIds []*entities.NbIdentity = nil
//	nodeIds := []*entities.NbIdentity{}
//	mockErr := &common.InternalError{Err: &net.OpError{Err:fmt.Errorf("connection error")}}
//	//readerMock.On("GetListNodebIds").Return(nilNodeIds, mockErr)
//	//readerMock.On("GetListNodebIds").Return(nodeIds, nil)
//
//	res, err := rnibDataService.GetListNodebIds()
//	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 2)
//	assert.True(t, strings.Contains(err.Error(),"connection failure", ))
//	assert.Equal(t, nodeIds, res)
//}

func TestGetE2TInstanceConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstanceOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	e2tInstance := &entities.E2TInstance{}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, nil)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstanceOkOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
