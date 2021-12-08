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
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
)

type RanLostConnectionHandler struct {
	ranDisconnectionManager managers.IRanDisconnectionManager
	logger                  *logger.Logger
}

func NewRanLostConnectionHandler(logger *logger.Logger, ranDisconnectionManager managers.IRanDisconnectionManager) RanLostConnectionHandler {
	return RanLostConnectionHandler{
		logger:                  logger,
		ranDisconnectionManager: ranDisconnectionManager,
	}
}
func (h RanLostConnectionHandler) Handle(request *models.NotificationRequest) {

	ranName := request.RanName

	h.logger.Warnf("#RanLostConnectionHandler.Handle - RAN name: %s - Received lost connection notification", ranName)

	_ = h.ranDisconnectionManager.DisconnectRan(ranName)
}
