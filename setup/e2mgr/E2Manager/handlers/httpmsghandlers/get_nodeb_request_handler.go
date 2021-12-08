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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
)

type GetNodebRequestHandler struct {
	rNibDataService services.RNibDataService
	logger          *logger.Logger
}

func NewGetNodebRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService) *GetNodebRequestHandler {
	return &GetNodebRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
	}
}

func (handler *GetNodebRequestHandler) Handle(request models.Request) (models.IResponse, error) {
	getNodebRequest := request.(models.GetNodebRequest)
	ranName:= getNodebRequest.RanName
	nodeb, err := handler.rNibDataService.GetNodeb(ranName)

	if err != nil {
		handler.logger.Errorf("#GetNodebRequestHandler.Handle - RAN name: %s - Error fetching RAN from rNib: %v",  ranName, err)
		return nil, rnibErrorToE2ManagerError(err)
	}

	return models.NewGetNodebResponse(nodeb), nil
}

func rnibErrorToE2ManagerError(err error) error {
	_, ok := err.(*common.ResourceNotFoundError)
	if ok {
		return e2managererrors.NewResourceNotFoundError()
	}
	return e2managererrors.NewRnibDbError()
}
