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

package httpmsghandlerprovider

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/handlers/httpmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type IncomingRequest string

const (
	ShutdownRequest        IncomingRequest = "Shutdown"
	ResetRequest           IncomingRequest = "Reset"
	X2SetupRequest         IncomingRequest = "X2SetupRequest"
	EndcSetupRequest       IncomingRequest = "EndcSetupRequest"
	GetNodebRequest        IncomingRequest = "GetNodebRequest"
	GetNodebIdListRequest  IncomingRequest = "GetNodebIdListRequest"
	GetE2TInstancesRequest IncomingRequest = "GetE2TInstancesRequest"
	UpdateGnbRequest       IncomingRequest = "UpdateGnbRequest"
)

type IncomingRequestHandlerProvider struct {
	requestMap map[IncomingRequest]httpmsghandlers.RequestHandler
	logger     *logger.Logger
}

func NewIncomingRequestHandlerProvider(logger *logger.Logger, rmrSender *rmrsender.RmrSender, config *configuration.Configuration, rNibDataService services.RNibDataService, ranSetupManager *managers.RanSetupManager, e2tInstancesManager managers.IE2TInstancesManager, e2tAssociationManager *managers.E2TAssociationManager, rmClient clients.IRoutingManagerClient) *IncomingRequestHandlerProvider {

	return &IncomingRequestHandlerProvider{
		requestMap: initRequestHandlerMap(logger, rmrSender, config, rNibDataService, ranSetupManager, e2tInstancesManager, e2tAssociationManager, rmClient),
		logger:     logger,
	}
}

func initRequestHandlerMap(logger *logger.Logger, rmrSender *rmrsender.RmrSender, config *configuration.Configuration, rNibDataService services.RNibDataService, ranSetupManager *managers.RanSetupManager, e2tInstancesManager managers.IE2TInstancesManager, e2tAssociationManager *managers.E2TAssociationManager, rmClient clients.IRoutingManagerClient) map[IncomingRequest]httpmsghandlers.RequestHandler {

	return map[IncomingRequest]httpmsghandlers.RequestHandler{
		ShutdownRequest:        httpmsghandlers.NewDeleteAllRequestHandler(logger, rmrSender, config, rNibDataService, e2tInstancesManager, rmClient),
		ResetRequest:           httpmsghandlers.NewX2ResetRequestHandler(logger, rmrSender, rNibDataService),
		X2SetupRequest:         httpmsghandlers.NewSetupRequestHandler(logger, rNibDataService, ranSetupManager, entities.E2ApplicationProtocol_X2_SETUP_REQUEST, e2tInstancesManager, e2tAssociationManager),
		EndcSetupRequest:       httpmsghandlers.NewSetupRequestHandler(logger, rNibDataService, ranSetupManager, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, e2tInstancesManager, e2tAssociationManager),
		GetNodebRequest:        httpmsghandlers.NewGetNodebRequestHandler(logger, rNibDataService),
		GetNodebIdListRequest:  httpmsghandlers.NewGetNodebIdListRequestHandler(logger, rNibDataService),
		GetE2TInstancesRequest: httpmsghandlers.NewGetE2TInstancesRequestHandler(logger, e2tInstancesManager),
		UpdateGnbRequest:       httpmsghandlers.NewUpdateGnbRequestHandler(logger, rNibDataService),
	}
}

func (provider IncomingRequestHandlerProvider) GetHandler(requestType IncomingRequest) (httpmsghandlers.RequestHandler, error) {
	handler, ok := provider.requestMap[requestType]

	if !ok {
		provider.logger.Errorf("#incoming_request_handler_provider.GetHandler - Cannot find handler for request type: %s", requestType)
		return nil, e2managererrors.NewInternalError()
	}

	return handler, nil
}
