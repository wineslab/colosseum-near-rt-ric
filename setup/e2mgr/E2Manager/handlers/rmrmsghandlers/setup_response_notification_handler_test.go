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


package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/converters"
	"e2mgr/e2managererrors"
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"unsafe"
)

const (
	RanName                           = "test"
	X2SetupResponsePackedPdu          = "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	EndcSetupResponsePackedPdu        = "202400808e00000100f600808640000200fc00090002f829504a952a0a00fd007200010c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000211e148033e4e5e4c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a00021a0044033e4e5e000000002c001e3f271f2e3d4ff0031e3f274400010000150400000a00020000"
	X2SetupFailureResponsePackedPdu   = "4006001a0000030005400200000016400100001140087821a00000008040"
	EndcSetupFailureResponsePackedPdu = "4024001a0000030005400200000016400100001140087821a00000008040"
)

type setupSuccessResponseTestCase struct {
	packedPdu            string
	setupResponseManager managers.ISetupResponseManager
	msgType              int
	saveNodebMockError   error
	sendMsgError         error
	statusChangeMbuf     *rmrCgo.MBuf
}

type setupFailureResponseTestCase struct {
	packedPdu            string
	setupResponseManager managers.ISetupResponseManager
	msgType              int
	saveNodebMockError   error
}

type setupResponseTestContext struct {
	logger                 *logger.Logger
	readerMock             *mocks.RnibReaderMock
	writerMock             *mocks.RnibWriterMock
	rnibDataService        services.RNibDataService
	setupResponseManager   managers.ISetupResponseManager
	ranStatusChangeManager managers.IRanStatusChangeManager
	rmrSender              *rmrsender.RmrSender
	rmrMessengerMock       *mocks.RmrMessengerMock
}

func NewSetupResponseTestContext(manager managers.ISetupResponseManager) *setupResponseTestContext {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	ranStatusChangeManager := managers.NewRanStatusChangeManager(logger, rmrSender)

	return &setupResponseTestContext{
		logger:                 logger,
		readerMock:             readerMock,
		writerMock:             writerMock,
		rnibDataService:        rnibDataService,
		setupResponseManager:   manager,
		ranStatusChangeManager: ranStatusChangeManager,
		rmrMessengerMock:       rmrMessengerMock,
		rmrSender:              rmrSender,
	}
}

func TestSetupResponseGetNodebFailure(t *testing.T) {
	notificationRequest := models.NotificationRequest{RanName: RanName}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.logger, testContext.rnibDataService, &managers.X2SetupResponseManager{}, testContext.ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP)
	testContext.readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewInternalError(errors.New("Error")))
	handler.Handle(&notificationRequest)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestSetupResponseInvalidConnectionStatus(t *testing.T) {
	ranName := "test"
	notificationRequest := models.NotificationRequest{RanName: ranName}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.logger, testContext.rnibDataService, &managers.X2SetupResponseManager{}, testContext.ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP)
	var rnibErr error
	testContext.readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}, rnibErr)
	handler.Handle(&notificationRequest)
	testContext.readerMock.AssertCalled(t, "GetNodeb", ranName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func executeHandleSetupSuccessResponse(t *testing.T, tc setupSuccessResponseTestCase) (*setupResponseTestContext, *entities.NodebInfo) {
	var payload []byte
	_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	notificationRequest := models.NotificationRequest{RanName: RanName, Payload: payload}
	testContext := NewSetupResponseTestContext(tc.setupResponseManager)

	handler := NewSetupResponseNotificationHandler(testContext.logger, testContext.rnibDataService, testContext.setupResponseManager, testContext.ranStatusChangeManager, tc.msgType)

	var rnibErr error

	nodebInfo := &entities.NodebInfo{
		ConnectionStatus:   entities.ConnectionStatus_CONNECTING,
		RanName:            RanName,
		Ip:                 "10.0.2.2",
		Port:               1231,
	}

	testContext.readerMock.On("GetNodeb", RanName).Return(nodebInfo, rnibErr)
	testContext.writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(tc.saveNodebMockError)
	testContext.rmrMessengerMock.On("SendMsg", tc.statusChangeMbuf, true).Return(&rmrCgo.MBuf{}, tc.sendMsgError)
	handler.Handle(&notificationRequest)

	return testContext, nodebInfo
}

func getRanConnectedMbuf(nodeType entities.Node_Type) *rmrCgo.MBuf {
	var xAction []byte
	resourceStatusPayload := models.NewResourceStatusPayload(nodeType, enums.RIC_TO_RAN)
	resourceStatusJson, _ := json.Marshal(resourceStatusPayload)
	var msgSrc unsafe.Pointer
	return rmrCgo.NewMBuf(rmrCgo.RAN_CONNECTED, len(resourceStatusJson), RanName, &resourceStatusJson, &xAction, msgSrc)
}

func executeHandleSetupFailureResponse(t *testing.T, tc setupFailureResponseTestCase) (*setupResponseTestContext, *entities.NodebInfo) {
	var payload []byte
	_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	notificationRequest := models.NotificationRequest{RanName: RanName, Payload: payload}
	testContext := NewSetupResponseTestContext(tc.setupResponseManager)

	handler := NewSetupResponseNotificationHandler(testContext.logger, testContext.rnibDataService, testContext.setupResponseManager, testContext.ranStatusChangeManager, tc.msgType)

	var rnibErr error

	nodebInfo := &entities.NodebInfo{
		ConnectionStatus:   entities.ConnectionStatus_CONNECTING,
		RanName:            RanName,
		Ip:                 "10.0.2.2",
		Port:               1231,
	}

	testContext.readerMock.On("GetNodeb", RanName).Return(nodebInfo, rnibErr)
	testContext.writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(tc.saveNodebMockError)
	handler.Handle(&notificationRequest)

	return testContext, nodebInfo
}

func TestX2SetupResponse(t *testing.T) {
	logger := initLog(t)
	var saveNodebMockError error
	var sendMsgError error
	tc := setupSuccessResponseTestCase{
		X2SetupResponsePackedPdu,
		managers.NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger)),
		rmrCgo.RIC_X2_SETUP_RESP,
		saveNodebMockError,
		sendMsgError,
		getRanConnectedMbuf(entities.Node_ENB),
	}

	testContext, nodebInfo := executeHandleSetupSuccessResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, entities.Node_ENB, nodebInfo.NodeType)

	assert.IsType(t, &entities.NodebInfo_Enb{}, nodebInfo.Configuration)
	i, _ := nodebInfo.Configuration.(*entities.NodebInfo_Enb)
	assert.NotNil(t, i.Enb)
	testContext.rmrMessengerMock.AssertCalled(t, "SendMsg", tc.statusChangeMbuf, true)
}

func TestX2SetupFailureResponse(t *testing.T) {
	logger := initLog(t)
	var saveNodebMockError error
	tc := setupFailureResponseTestCase{
		X2SetupFailureResponsePackedPdu,
		managers.NewX2SetupFailureResponseManager(converters.NewX2SetupFailureResponseConverter(logger)),
		rmrCgo.RIC_X2_SETUP_FAILURE,
		saveNodebMockError,
	}

	testContext, nodebInfo := executeHandleSetupFailureResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, entities.Failure_X2_SETUP_FAILURE, nodebInfo.FailureType)
	assert.NotNil(t, nodebInfo.SetupFailure)
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestEndcSetupResponse(t *testing.T) {
	logger := initLog(t)
	var saveNodebMockError error
	var sendMsgError error
	tc := setupSuccessResponseTestCase{
		EndcSetupResponsePackedPdu,
		managers.NewEndcSetupResponseManager(converters.NewEndcSetupResponseConverter(logger)),
		rmrCgo.RIC_ENDC_X2_SETUP_RESP,
		saveNodebMockError,
		sendMsgError,
		getRanConnectedMbuf(entities.Node_GNB),
	}

	testContext, nodebInfo := executeHandleSetupSuccessResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, entities.Node_GNB, nodebInfo.NodeType)
	assert.IsType(t, &entities.NodebInfo_Gnb{}, nodebInfo.Configuration)

	i, _ := nodebInfo.Configuration.(*entities.NodebInfo_Gnb)
	assert.NotNil(t, i.Gnb)
	testContext.rmrMessengerMock.AssertCalled(t, "SendMsg", tc.statusChangeMbuf, true)
}

func TestEndcSetupFailureResponse(t *testing.T) {
	logger := initLog(t)
	var saveNodebMockError error
	tc := setupFailureResponseTestCase{
		EndcSetupFailureResponsePackedPdu,
		managers.NewEndcSetupFailureResponseManager(converters.NewEndcSetupFailureResponseConverter(logger)),
		rmrCgo.RIC_ENDC_X2_SETUP_FAILURE,
		saveNodebMockError,
	}

	testContext, nodebInfo := executeHandleSetupFailureResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, entities.Failure_ENDC_X2_SETUP_FAILURE, nodebInfo.FailureType)
	assert.NotNil(t, nodebInfo.SetupFailure)
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestSetupResponseInvalidPayload(t *testing.T) {
	logger := initLog(t)
	ranName := "test"
	notificationRequest := models.NotificationRequest{RanName: ranName, Payload: []byte("123")}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.logger, testContext.rnibDataService, managers.NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger)), testContext.ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP)
	var rnibErr error
	testContext.readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING}, rnibErr)
	handler.Handle(&notificationRequest)
	testContext.readerMock.AssertCalled(t, "GetNodeb", ranName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestSetupResponseSaveNodebFailure(t *testing.T) {
	logger := initLog(t)
	saveNodebMockError := common.NewInternalError(errors.New("Error"))
	var sendMsgError error
	tc := setupSuccessResponseTestCase{
		X2SetupResponsePackedPdu,
		managers.NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger)),
		rmrCgo.RIC_X2_SETUP_RESP,
		saveNodebMockError,
		sendMsgError,
		getRanConnectedMbuf(entities.Node_ENB),
	}

	testContext, nodebInfo := executeHandleSetupSuccessResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	testContext.rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestSetupResponseStatusChangeSendFailure(t *testing.T) {
	logger := initLog(t)
	var saveNodebMockError error
	sendMsgError := e2managererrors.NewRmrError()
	tc := setupSuccessResponseTestCase{
		X2SetupResponsePackedPdu,
		managers.NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger)),
		rmrCgo.RIC_X2_SETUP_RESP,
		saveNodebMockError,
		sendMsgError,
		getRanConnectedMbuf(entities.Node_ENB),
	}

	testContext, nodebInfo := executeHandleSetupSuccessResponse(t, tc)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, entities.Node_ENB, nodebInfo.NodeType)

	assert.IsType(t, &entities.NodebInfo_Enb{}, nodebInfo.Configuration)
	i, _ := nodebInfo.Configuration.(*entities.NodebInfo_Enb)
	assert.NotNil(t, i.Enb)
	testContext.rmrMessengerMock.AssertCalled(t, "SendMsg", tc.statusChangeMbuf, true)
}
