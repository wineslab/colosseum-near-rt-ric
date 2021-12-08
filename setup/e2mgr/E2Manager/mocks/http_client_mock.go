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

package mocks

import (
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
)

type HttpClientMock struct {
	mock.Mock
}

func (c *HttpClientMock) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := c.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *HttpClientMock) Delete(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := c.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}
