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

package clients

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"encoding/json"
	"net/http"
)

const (
	AddE2TInstanceApiSuffix            = "e2t"
	AssociateRanToE2TInstanceApiSuffix = "associate-ran-to-e2t"
	DissociateRanE2TInstanceApiSuffix  = "dissociate-ran"
	DeleteE2TInstanceApiSuffix         = "e2t"
)

type RoutingManagerClient struct {
	logger     *logger.Logger
	config     *configuration.Configuration
	httpClient IHttpClient
}

type IRoutingManagerClient interface {
	AddE2TInstance(e2tAddress string) error
	AssociateRanToE2TInstance(e2tAddress string, ranName string) error
	DissociateRanE2TInstance(e2tAddress string, ranName string) error
	DissociateAllRans(e2tAddresses []string) error
	DeleteE2TInstance(e2tAddress string, ransToBeDissociated []string) error
}

func NewRoutingManagerClient(logger *logger.Logger, config *configuration.Configuration, httpClient IHttpClient) *RoutingManagerClient {
	return &RoutingManagerClient{
		logger:     logger,
		config:     config,
		httpClient: httpClient,
	}
}

func (c *RoutingManagerClient) AddE2TInstance(e2tAddress string) error {

	data := models.NewRoutingManagerE2TData(e2tAddress)
	url := c.config.RoutingManager.BaseUrl + AddE2TInstanceApiSuffix

	return c.PostMessage(url, data)
}

func (c *RoutingManagerClient) AssociateRanToE2TInstance(e2tAddress string, ranName string) error {

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(e2tAddress, ranName)}
	url := c.config.RoutingManager.BaseUrl + AssociateRanToE2TInstanceApiSuffix

	return c.PostMessage(url, data)
}

func (c *RoutingManagerClient) DissociateRanE2TInstance(e2tAddress string, ranName string) error {

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(e2tAddress, ranName)}
	url := c.config.RoutingManager.BaseUrl + DissociateRanE2TInstanceApiSuffix

	return c.PostMessage(url, data)
}

func (c *RoutingManagerClient) DissociateAllRans(e2tAddresses []string) error {

	data := mapE2TAddressesToE2DataList(e2tAddresses)
	url := c.config.RoutingManager.BaseUrl + DissociateRanE2TInstanceApiSuffix

	return c.PostMessage(url, data)
}

func (c *RoutingManagerClient) DeleteE2TInstance(e2tAddress string, ransTobeDissociated []string) error {
	data := models.NewRoutingManagerDeleteRequestModel(e2tAddress, ransTobeDissociated, nil)
	url := c.config.RoutingManager.BaseUrl + DeleteE2TInstanceApiSuffix
	return c.DeleteMessage(url, data)
}

func (c *RoutingManagerClient) sendMessage(method string, url string, data interface{}) error {
	marshaled, err := json.Marshal(data)

	if err != nil {
		return e2managererrors.NewRoutingManagerError()
	}

	body := bytes.NewBuffer(marshaled)
	c.logger.Infof("[E2 Manager -> Routing Manager] #RoutingManagerClient.sendMessage - %s url: %s, request body: %+v", method, url, body)

	var resp *http.Response

	if method == http.MethodPost {
		resp, err = c.httpClient.Post(url, "application/json", body)
	} else if method == http.MethodDelete {
		resp, err = c.httpClient.Delete(url, "application/json", body)
	}

	if err != nil {
		c.logger.Errorf("#RoutingManagerClient.sendMessage - failed sending request. error: %s", err)
		return e2managererrors.NewRoutingManagerError()
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		c.logger.Infof("[Routing Manager -> E2 Manager] #RoutingManagerClient.sendMessage - success. http status code: %d", resp.StatusCode)
		return nil
	}

	c.logger.Errorf("[Routing Manager -> E2 Manager] #RoutingManagerClient.sendMessage - failure. http status code: %d", resp.StatusCode)
	return e2managererrors.NewRoutingManagerError()
}

func (c *RoutingManagerClient) DeleteMessage(url string, data interface{}) error {
	return c.sendMessage(http.MethodDelete, url, data)
}

func (c *RoutingManagerClient) PostMessage(url string, data interface{}) error {
	return c.sendMessage(http.MethodPost, url, data)
}

func mapE2TAddressesToE2DataList(e2tAddresses []string) models.RoutingManagerE2TDataList {
	e2tDataList := make(models.RoutingManagerE2TDataList, len(e2tAddresses))

	for i, v := range e2tAddresses {
		e2tDataList[i] = models.NewRoutingManagerE2TData(v)
	}

	return e2tDataList
}

func convertE2TToRansAssociationsMapToE2TDataList(e2tToRansAssociations map[string][]string) models.RoutingManagerE2TDataList {
	e2tDataList := make(models.RoutingManagerE2TDataList, len(e2tToRansAssociations))
	i := 0
	for k, v := range e2tToRansAssociations {
		e2tDataList[i] = models.NewRoutingManagerE2TData(k, v...)
		i++
	}

	return e2tDataList
}
