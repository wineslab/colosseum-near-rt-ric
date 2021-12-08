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


package rmrreceiver

import (
	"e2mgr/logger"
	"e2mgr/managers/notificationmanager"
	"e2mgr/rmrCgo"
)

type RmrReceiver struct {
	logger    *logger.Logger
	nManager  *notificationmanager.NotificationManager
	messenger rmrCgo.RmrMessenger
}

func NewRmrReceiver(logger *logger.Logger, messenger rmrCgo.RmrMessenger, nManager *notificationmanager.NotificationManager) *RmrReceiver {
	return &RmrReceiver{
		logger:    logger,
		nManager:  nManager,
		messenger: messenger,
	}
}

func (r *RmrReceiver) ListenAndHandle() {

	for {
		mbuf, err := r.messenger.RecvMsg()

		if err != nil {
			r.logger.Errorf("#RmrReceiver.ListenAndHandle - error: %s", err)
			continue
		}

		r.logger.Debugf("#RmrReceiver.ListenAndHandle - Going to handle received message: %#v\n", mbuf)

		// TODO: go routine?
		_ = r.nManager.HandleMessage(mbuf)
	}
}