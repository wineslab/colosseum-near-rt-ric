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

package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services/rmrsender"
	"time"
)

type E2TKeepAliveWorker struct {
	logger              *logger.Logger
	e2tShutdownManager  IE2TShutdownManager
	e2TInstancesManager IE2TInstancesManager
	rmrSender           *rmrsender.RmrSender
	config              *configuration.Configuration
}

func NewE2TKeepAliveWorker(logger *logger.Logger, rmrSender *rmrsender.RmrSender, e2TInstancesManager IE2TInstancesManager, e2tShutdownManager IE2TShutdownManager, config *configuration.Configuration) E2TKeepAliveWorker {
	return E2TKeepAliveWorker{
		logger:              logger,
		e2tShutdownManager:  e2tShutdownManager,
		e2TInstancesManager: e2TInstancesManager,
		rmrSender:           rmrSender,
		config:              config,
	}
}

func (h E2TKeepAliveWorker) Execute() {

	h.logger.Infof("#E2TKeepAliveWorker.Execute - keep alive started")

	ticker := time.NewTicker(time.Duration(h.config.KeepAliveDelayMs) * time.Millisecond)

	for _ = range ticker.C {

		h.SendKeepAliveRequest()
		h.E2TKeepAliveExpired()
	}
}

func (h E2TKeepAliveWorker) E2TKeepAliveExpired() {

	e2tInstances, err := h.e2TInstancesManager.GetE2TInstancesNoLogs()

	if err != nil || len(e2tInstances) == 0 {
		return
	}

	for _, e2tInstance := range e2tInstances {

		delta := int64(time.Now().UnixNano()) - e2tInstance.KeepAliveTimestamp
		timestampNanosec := int64(time.Duration(h.config.KeepAliveResponseTimeoutMs) * time.Millisecond)

		if delta > timestampNanosec {

			h.logger.Warnf("#E2TKeepAliveWorker.E2TKeepAliveExpired - e2t address: %s time expired, shutdown e2 instance", e2tInstance.Address)

			h.e2tShutdownManager.Shutdown(e2tInstance)
		}
	}
}

func (h E2TKeepAliveWorker) SendKeepAliveRequest() {

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.E2_TERM_KEEP_ALIVE_REQ}
	h.rmrSender.SendWithoutLogs(&rmrMessage)
}
