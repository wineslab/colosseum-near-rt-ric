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

package httpmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
)

type GetE2TInstancesRequestHandler struct {
	e2tInstancesManager managers.IE2TInstancesManager
	logger              *logger.Logger
}

func NewGetE2TInstancesRequestHandler(logger *logger.Logger, e2tInstancesManager managers.IE2TInstancesManager) *GetE2TInstancesRequestHandler {
	return &GetE2TInstancesRequestHandler{
		logger:              logger,
		e2tInstancesManager: e2tInstancesManager,
	}
}

func (h *GetE2TInstancesRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	e2tInstances, err := h.e2tInstancesManager.GetE2TInstances()

	if err != nil {
		h.logger.Errorf("#GetE2TInstancesRequestHandler.Handle - Error fetching E2T instances from rNib: %s", err)
		return nil, err
	}

	mapped := make([]*models.E2TInstanceResponseModel, len(e2tInstances))

	for i, v := range e2tInstances {
		mapped[i] = models.NewE2TInstanceResponseModel(v.Address, v.AssociatedRanList)
	}

	return models.E2TInstancesResponse(mapped), nil
}
