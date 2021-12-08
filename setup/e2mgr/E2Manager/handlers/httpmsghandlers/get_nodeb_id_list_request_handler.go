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


package httpmsghandlers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
)

type GetNodebIdListRequestHandler struct {
	rNibDataService services.RNibDataService
	logger          *logger.Logger
}

func NewGetNodebIdListRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService) *GetNodebIdListRequestHandler {
	return &GetNodebIdListRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
	}
}

func (handler *GetNodebIdListRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	nodebIdList, err := handler.rNibDataService.GetListNodebIds()

	if err != nil {
		handler.logger.Errorf("#GetNodebIdListRequestHandler.Handle - Error fetching Nodeb Identity list from rNib: %v", err)
		return nil, e2managererrors.NewRnibDbError()
	}

	return models.NewGetNodebIdListResponse(nodebIdList), nil
}
