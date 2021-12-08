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

package httpserver

import (
	"e2mgr/controllers"
	"e2mgr/logger"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Run(log *logger.Logger, port int, rootController controllers.IRootController, nodebController controllers.INodebController, e2tController controllers.IE2TController) error {

	router := mux.NewRouter();
	initializeRoutes(router, rootController, nodebController, e2tController)

	addr := fmt.Sprintf(":%d", port)

	err := http.ListenAndServe(addr, router)

	log.Errorf("#http_server.Run - Fail initiating HTTP server. Error: %v", err)
	return err
}

func initializeRoutes(router *mux.Router, rootController controllers.IRootController, nodebController controllers.INodebController, e2tController controllers.IE2TController) {
	r := router.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/health", rootController.HandleHealthCheckRequest).Methods(http.MethodGet)

	rr := r.PathPrefix("/nodeb").Subrouter()
	rr.HandleFunc("/ids", nodebController.GetNodebIdList).Methods(http.MethodGet)
	rr.HandleFunc("/{ranName}", nodebController.GetNodeb).Methods(http.MethodGet)
	rr.HandleFunc("/{ranName}/update", nodebController.UpdateGnb).Methods(http.MethodPut)
	rr.HandleFunc("/shutdown", nodebController.Shutdown).Methods(http.MethodPut)
	rrr := r.PathPrefix("/e2t").Subrouter()
	rrr.HandleFunc("/list", e2tController.GetE2TInstances).Methods(http.MethodGet)
}
