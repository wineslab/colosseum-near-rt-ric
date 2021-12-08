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

/*import (
	"e2mgr/models"
	"fmt"
	"time"
)

const (
	FullUperPdu  string = "004c07080004001980da0100075bde017c148003d5a8205000017c180003d5a875555403331420000012883a0003547400cd20002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be0c4001ead4016e007ab50100002f8320067ab5005b8c1ead5070190c00001d637805f220000f56a081400005f020000f56a1d555400ccc508002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be044001ead4016e007ab50100002f8120067ab5005b8c1ead5070190c00000"
	FullAperPdu  string = "" // TODO: populate and use it
	BasicUperPdu string = "004898000400190d0000074200017c148003d5a80000"
	BasicAperPdu string = "" // TODO: populate and use it
	GarbagePdu   string = "12312312"
)

func createNotificationRequest(ranName string, transactionId []byte, packedPdu string) (*models.NotificationRequest, error) {
	var packedByteSlice []byte

	_, err := fmt.Sscanf(packedPdu, "%x", &packedByteSlice)

	if err != nil {
		return nil, err
	}

	return models.NewNotificationRequest(ranName, packedByteSlice, time.Now(), transactionId, nil), nil
}

func createNotificationRequestAndHandle(ranName string, transactionId []byte, loadInformationHandler EnbLoadInformationNotificationHandler, pdu string) error {
	notificationRequest, err := createNotificationRequest(ranName, transactionId, pdu)

	if err != nil {
		return err
	}

	loadInformationHandler.Handle(notificationRequest)
	return nil
}*/

//func TestLoadInformationHandlerSuccess(t *testing.T) {
//	log, err := logger.InitLogger(logger.InfoLevel)
//	if err != nil {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerSuccess - failed to initialize logger, error: %v", err)
//	}
//
//	inventoryName := "testRan"
//
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//
//	var rnibErr error
//	writerMock.On("SaveRanLoadInformation",inventoryName, mock.Anything).Return(rnibErr)
//
//	loadInformationHandler := NewEnbLoadInformationNotificationHandler(rnibWriterProvider)
//
//	var packedExampleByteSlice []byte
//	_, err = fmt.Sscanf(FullUperPdu, "%x", &packedExampleByteSlice)
//	notificationRequest := models.NewNotificationRequest(inventoryName, packedExampleByteSlice, time.Now(), " 881828026419")
//	loadInformationHandler.Handle(log, notificationRequest)
//
//	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 1)
//}
//
//func TestLoadInformationHandlerPayloadFailure(t *testing.T) {
//	log, err := logger.InitLogger(logger.InfoLevel)
//	if err != nil {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerPayloadFailure - failed to initialize logger, error: %v", err)
//	}
//
//	inventoryName := "testRan"
//
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//
//	var rnibErr error
//	writerMock.On("SaveRanLoadInformation",inventoryName, mock.Anything).Return(rnibErr)
//
//	loadInformationHandler := NewEnbLoadInformationNotificationHandler(rnibWriterProvider)
//
//	var packedExampleByteSlice []byte
//	_, err = fmt.Sscanf(GarbagePdu, "%x", &packedExampleByteSlice)
//	notificationRequest := models.NewNotificationRequest(inventoryName, packedExampleByteSlice, time.Now(), " 881828026419")
//	loadInformationHandler.Handle(log, notificationRequest)
//
//	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 0)
//}

// Integration test
//func TestLoadInformationHandlerOverrideSuccess(t *testing.T) {
//	log, err := logger.InitLogger(logger.InfoLevel)
//	if err != nil {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerOverrideSuccess - failed to initialize logger, error: %s", err)
//	}
//
//	rNibWriter.Init("e2Manager", 1)
//	defer rNibWriter.Close()
//	reader.Init("e2Manager", 1)
//	defer reader.Close()
//	loadInformationHandler := NewEnbLoadInformationNotificationHandler(rNibWriter.GetRNibWriter)
//
//	err = createNotificationRequestAndHandle("ranName", " 881828026419", loadInformationHandler, FullUperPdu)
//
//	if err != nil {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerOverrideSuccess - failed creating NotificationRequest, error: %v", err)
//	}
//
//	err = createNotificationRequestAndHandle("ranName", " 881828026419", loadInformationHandler, BasicUperPdu)
//
//	if err != nil {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerOverrideSuccess - failed creating NotificationRequest, error: %v", err)
//	}
//
//	ranLoadInformation, rnibErr := reader.GetRNibReader().GetRanLoadInformation("ranName")
//
//	if (rnibErr != nil) {
//		t.Errorf("#setup_request_handler_test.TestLoadInformationHandlerOverrideSuccess - RNIB error: %v", err)
//	}
//
//	assert.Len(t, ranLoadInformation.CellLoadInfos, 1)
//}
