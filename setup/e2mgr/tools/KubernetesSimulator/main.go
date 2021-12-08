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


package main

import (
	"fmt"
	"kubsimulator/configuration"
	"kubsimulator/go"
	"log"
	"net/http"
)

func main() {
	config := configuration.ParseConfiguration()
	port := config.Http.Port

	log.Printf("Server started on port %d", port)

	router := kubernetes.NewRouter()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
