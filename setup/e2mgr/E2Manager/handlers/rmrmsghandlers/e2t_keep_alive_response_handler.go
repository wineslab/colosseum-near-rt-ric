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
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package rmrmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/services"
	"encoding/json"
)

type E2TKeepAliveResponseHandler struct {
	logger              *logger.Logger
	rnibDataService     services.RNibDataService
	e2TInstancesManager managers.IE2TInstancesManager
}

func NewE2TKeepAliveResponseHandler(logger *logger.Logger, rnibDataService services.RNibDataService, e2TInstancesManager managers.IE2TInstancesManager) E2TKeepAliveResponseHandler {
	return E2TKeepAliveResponseHandler{
		logger:              logger,
		rnibDataService:     rnibDataService,
		e2TInstancesManager: e2TInstancesManager,
	}
}

func (h E2TKeepAliveResponseHandler) Handle(request *models.NotificationRequest) {
	unmarshalledPayload := models.E2TKeepAlivePayload{}
	err := json.Unmarshal(request.Payload, &unmarshalledPayload)

	if err != nil {
		h.logger.Errorf("#E2TKeepAliveResponseHandler.Handle - Error unmarshaling RMR request payload: %v", err)
		return
	}

	_ = h.e2TInstancesManager.ResetKeepAliveTimestamp(unmarshalledPayload.Address)
}
