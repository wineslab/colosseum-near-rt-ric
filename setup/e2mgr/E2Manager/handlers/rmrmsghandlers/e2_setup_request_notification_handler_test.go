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
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	prefix                   = "10.0.2.15:9999|"
	e2tInstanceFullAddress   = "10.0.2.15:9999"
	nodebRanName             = "gnb:310-410-b5c67788"
	GnbSetupRequestXmlPath   = "../../tests/resources/setupRequest_gnb.xml"
	GnbWithoutFunctionsSetupRequestXmlPath   = "../../tests/resources/setupRequest_gnb_without_functions.xml"
	EnGnbSetupRequestXmlPath = "../../tests/resources/setupRequest_en-gNB.xml"
	NgEnbSetupRequestXmlPath = "../../tests/resources/setupRequest_ng-eNB.xml"
	EnbSetupRequestXmlPath   = "../../tests/resources/setupRequest_enb.xml"
)

func readXmlFile(t *testing.T, xmlPath string) []byte {
	path, err := filepath.Abs(xmlPath)
	if err != nil {
		t.Fatal(err)
	}
	xmlAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return xmlAsBytes
}

func TestParseGnbSetupRequest_Success(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Equal(t, "131014", request.GetPlmnId())
	assert.Equal(t, "10011001101010101011", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseEnGnbSetupRequest_Success(t *testing.T) {
	enGnbXml := readXmlFile(t, EnGnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, enGnbXml...))
	assert.Equal(t, "131014", request.GetPlmnId())
	assert.Equal(t, "11000101110001101100011111111000", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseNgEnbSetupRequest_Success(t *testing.T) {
	ngEnbXml := readXmlFile(t, NgEnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, ngEnbXml...))
	assert.Equal(t, "131014", request.GetPlmnId())
	assert.Equal(t, "101010101010101010", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseEnbSetupRequest_Success(t *testing.T) {
	enbXml := readXmlFile(t, EnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, enbXml...))
	assert.Equal(t, "6359AB", request.GetPlmnId())
	assert.Equal(t, "101010101010101010", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseSetupRequest_PipFailure(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte("10.0.2.15:9999")
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
}

func TestParseSetupRequest_UnmarshalFailure(t *testing.T) {
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, 1, 2, 3))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: 31302e302e322e31353a393939397c010203")
}

func TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertNewNodebSuccessCalls(readerMock, t, e2tInstancesManagerMock, writerMock, routingManagerClientMock, rmrMessengerMock)
}

func TestE2SetupRequestNotificationHandler_HandleNewGnbWithoutFunctionsSuccess(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbWithoutFunctionsSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertNewNodebSuccessCalls(readerMock, t, e2tInstancesManagerMock, writerMock, routingManagerClientMock, rmrMessengerMock)
}

func TestE2SetupRequestNotificationHandler_HandleNewEnGnbSuccess(t *testing.T) {
	xmlEnGnb := readXmlFile(t, EnGnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlEnGnb...)}
	handler.Handle(notificationRequest)
	assertNewNodebSuccessCalls(readerMock, t, e2tInstancesManagerMock, writerMock, routingManagerClientMock, rmrMessengerMock)
}

func TestE2SetupRequestNotificationHandler_HandleNewNgEnbSuccess(t *testing.T) {
	xmlNgEnb := readXmlFile(t, NgEnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlNgEnb...)}
	handler.Handle(notificationRequest)
	assertNewNodebSuccessCalls(readerMock, t, e2tInstancesManagerMock, writerMock, routingManagerClientMock, rmrMessengerMock)
}

func TestE2SetupRequestNotificationHandler_HandleExistingGnbSuccess(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)

	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb = &entities.NodebInfo{
		RanName:                      nodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertExistingNodebSuccessCalls(readerMock, t, e2tInstancesManagerMock, writerMock, routingManagerClientMock, rmrMessengerMock)
}

func TestE2SetupRequestNotificationHandler_HandleParseError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)

	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	prefBytes := []byte("invalid_prefix")
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleUnmarshalError(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, "xmlGnb"...)}
	handler.Handle(notificationRequest)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleGetE2TInstanceError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance *entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, common.NewResourceNotFoundError("Not found"))
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleGetNodebError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)

	handler, readerMock, writerMock, routingManagerClientMock, e2tInstancesManagerMock, rmrMessengerMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewInternalError(errors.New("some error")))
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleAssociationError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)

	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(errors.New("association error"))
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("WhSendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)

	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertCalled(t, "WhSendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_ConvertTo20BitStringError(t *testing.T) {
	xmlEnGnb := readXmlFile(t, EnGnbSetupRequestXmlPath)
	logger := tests.InitLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, GlobalRicId: struct {
		PlmnId      string
		RicNearRtId string
	}{PlmnId: "131014", RicNearRtId: "10011001101010101011"}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock)
	handler := NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManagerMock, rmrSender, rnibDataService, e2tAssociationManager)

	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlEnGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleExistingGnbInvalidStatusError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, routingManagerClientMock, e2tInstancesManagerMock, rmrMessengerMock := initMocks(t)
	var gnb = &entities.NodebInfo{RanName: nodebRanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, nil)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func initMocks(t *testing.T) (E2SetupRequestNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, GlobalRicId: struct {
		PlmnId      string
		RicNearRtId string
	}{PlmnId: "131014", RicNearRtId: "556670"}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock)
	handler := NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManagerMock, rmrSender, rnibDataService, e2tAssociationManager)
	return handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock
}

func assertNewNodebSuccessCalls(readerMock *mocks.RnibReaderMock, t *testing.T, e2tInstancesManagerMock *mocks.E2TInstancesManagerMock, writerMock *mocks.RnibWriterMock, routingManagerClientMock *mocks.RoutingManagerClientMock, rmrMessengerMock *mocks.RmrMessengerMock) {
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func assertExistingNodebSuccessCalls(readerMock *mocks.RnibReaderMock, t *testing.T, e2tInstancesManagerMock *mocks.E2TInstancesManagerMock, writerMock *mocks.RnibWriterMock, routingManagerClientMock *mocks.RoutingManagerClientMock, rmrMessengerMock *mocks.RmrMessengerMock) {
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mock.Anything, mock.Anything)
}
