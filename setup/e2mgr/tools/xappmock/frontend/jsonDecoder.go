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

package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"xappmock/models"
)

func DecodeJsonCommand(data []byte) (*models.JsonCommand, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var cmd models.JsonCommand
	if err := dec.Decode(&cmd); err != nil && err != io.EOF {
		return nil, errors.New(err.Error())
	}

	return &cmd, nil
}

func JsonCommandsDecoder(data []byte, processor func(models.JsonCommand) error) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	for {
		var commands []models.JsonCommand
		if err := dec.Decode(&commands); err == io.EOF {
			break
		} else if err != nil {
			return errors.New(err.Error())
		}
		for i, cmd := range commands {
			if err := processor(cmd); err != nil {
				return errors.New(fmt.Sprintf("processing error at #%d, %s", i, err))
			}
		}
	}
	return nil
}
