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


package notificationmanager

import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rmrCgo"
	"time"
)

type NotificationManager struct {
	logger                      *logger.Logger
	notificationHandlerProvider *rmrmsghandlerprovider.NotificationHandlerProvider
}

func NewNotificationManager(logger *logger.Logger, notificationHandlerProvider *rmrmsghandlerprovider.NotificationHandlerProvider) *NotificationManager {
	return &NotificationManager{
		logger:                      logger,
		notificationHandlerProvider: notificationHandlerProvider,
	}
}

func (m NotificationManager) HandleMessage(mbuf *rmrCgo.MBuf) error {

	notificationHandler, err := m.notificationHandlerProvider.GetNotificationHandler(mbuf.MType)

	if err != nil {
		m.logger.Errorf("#NotificationManager.HandleMessage - Error: %s", err)
		return err
	}

	notificationRequest := models.NewNotificationRequest(mbuf.Meid, *mbuf.Payload, time.Now(), *mbuf.XAction, mbuf.GetMsgSrc())
	go notificationHandler.Handle(notificationRequest)
	return nil
}
