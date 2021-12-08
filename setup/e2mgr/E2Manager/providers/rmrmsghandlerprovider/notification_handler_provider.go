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

package rmrmsghandlerprovider

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/converters"
	"e2mgr/handlers/rmrmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"fmt"
)

type NotificationHandlerProvider struct {
	notificationHandlers map[int]rmrmsghandlers.NotificationHandler
}

func NewNotificationHandlerProvider() *NotificationHandlerProvider {
	return &NotificationHandlerProvider{
		notificationHandlers: map[int]rmrmsghandlers.NotificationHandler{},
	}
}

// TODO: check whether it has been initialized
func (provider NotificationHandlerProvider) GetNotificationHandler(messageType int) (rmrmsghandlers.NotificationHandler, error) {
	handler, ok := provider.notificationHandlers[messageType]

	if !ok {
		return nil, fmt.Errorf("notification handler not found for message %d", messageType)
	}

	return handler, nil
}

func (provider *NotificationHandlerProvider) Register(msgType int, handler rmrmsghandlers.NotificationHandler) {
	provider.notificationHandlers[msgType] = handler
}

func (provider *NotificationHandlerProvider) Init(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, rmrSender *rmrsender.RmrSender, ranSetupManager *managers.RanSetupManager, e2tInstancesManager managers.IE2TInstancesManager, routingManagerClient clients.IRoutingManagerClient, e2tAssociationManager *managers.E2TAssociationManager) {

	// Init converters
	x2SetupResponseConverter := converters.NewX2SetupResponseConverter(logger)
	x2SetupFailureResponseConverter := converters.NewX2SetupFailureResponseConverter(logger)
	endcSetupResponseConverter := converters.NewEndcSetupResponseConverter(logger)
	endcSetupFailureResponseConverter := converters.NewEndcSetupFailureResponseConverter(logger)
	//enbLoadInformationExtractor := converters.NewEnbLoadInformationExtractor(logger)
	x2ResetResponseExtractor := converters.NewX2ResetResponseExtractor(logger)

	// Init managers
	ranReconnectionManager := managers.NewRanDisconnectionManager(logger, config, rnibDataService, e2tAssociationManager)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(logger, rmrSender)
	x2SetupResponseManager := managers.NewX2SetupResponseManager(x2SetupResponseConverter)
	x2SetupFailureResponseManager := managers.NewX2SetupFailureResponseManager(x2SetupFailureResponseConverter)
	endcSetupResponseManager := managers.NewEndcSetupResponseManager(endcSetupResponseConverter)
	endcSetupFailureResponseManager := managers.NewEndcSetupFailureResponseManager(endcSetupFailureResponseConverter)

	// Init handlers
	x2SetupResponseHandler := rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupResponseManager, ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP)
	x2SetupFailureResponseHandler := rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupFailureResponseManager, nil, rmrCgo.RIC_X2_SETUP_FAILURE)
	endcSetupResponseHandler := rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, endcSetupResponseManager, ranStatusChangeManager, rmrCgo.RIC_ENDC_X2_SETUP_RESP)
	endcSetupFailureResponseHandler := rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, endcSetupFailureResponseManager, nil, rmrCgo.RIC_ENDC_X2_SETUP_FAILURE)
	ranLostConnectionHandler := rmrmsghandlers.NewRanLostConnectionHandler(logger, ranReconnectionManager)
	//enbLoadInformationNotificationHandler := rmrmsghandlers.NewEnbLoadInformationNotificationHandler(logger, rnibDataService, enbLoadInformationExtractor)
	x2EnbConfigurationUpdateHandler := rmrmsghandlers.NewX2EnbConfigurationUpdateHandler(logger, rmrSender)
	endcConfigurationUpdateHandler := rmrmsghandlers.NewEndcConfigurationUpdateHandler(logger, rmrSender)
	x2ResetResponseHandler := rmrmsghandlers.NewX2ResetResponseHandler(logger, rnibDataService, ranStatusChangeManager, x2ResetResponseExtractor)
	x2ResetRequestNotificationHandler := rmrmsghandlers.NewX2ResetRequestNotificationHandler(logger, rnibDataService, ranStatusChangeManager, rmrSender)
	e2TermInitNotificationHandler := rmrmsghandlers.NewE2TermInitNotificationHandler(logger, ranReconnectionManager, e2tInstancesManager, routingManagerClient)
	e2TKeepAliveResponseHandler := rmrmsghandlers.NewE2TKeepAliveResponseHandler(logger, rnibDataService, e2tInstancesManager)
	e2SetupRequestNotificationHandler := rmrmsghandlers.NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManager, rmrSender, rnibDataService, e2tAssociationManager)

	provider.Register(rmrCgo.RIC_X2_SETUP_RESP, x2SetupResponseHandler)
	provider.Register(rmrCgo.RIC_X2_SETUP_FAILURE, x2SetupFailureResponseHandler)
	provider.Register(rmrCgo.RIC_ENDC_X2_SETUP_RESP, endcSetupResponseHandler)
	provider.Register(rmrCgo.RIC_ENDC_X2_SETUP_FAILURE, endcSetupFailureResponseHandler)
	provider.Register(rmrCgo.RIC_SCTP_CONNECTION_FAILURE, ranLostConnectionHandler)
	//provider.Register(rmrCgo.RIC_ENB_LOAD_INFORMATION, enbLoadInformationNotificationHandler)
	provider.Register(rmrCgo.RIC_ENB_CONF_UPDATE, x2EnbConfigurationUpdateHandler)
	provider.Register(rmrCgo.RIC_ENDC_CONF_UPDATE, endcConfigurationUpdateHandler)
	provider.Register(rmrCgo.RIC_X2_RESET_RESP, x2ResetResponseHandler)
	provider.Register(rmrCgo.RIC_X2_RESET, x2ResetRequestNotificationHandler)
	provider.Register(rmrCgo.RIC_E2_TERM_INIT, e2TermInitNotificationHandler)
	provider.Register(rmrCgo.E2_TERM_KEEP_ALIVE_RESP, e2TKeepAliveResponseHandler)
	provider.Register(rmrCgo.RIC_E2_SETUP_REQ, e2SetupRequestNotificationHandler)
}
