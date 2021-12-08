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
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ProcessConfigurationFile(resourcesFolder, inputFolder, suffix string, processor func(data []byte) error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.New(err.Error())
	}
	inputDir := filepath.Join(cwd, resourcesFolder, inputFolder)

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return errors.New(err.Error())
	}

	for _, file := range files {
		if file.Mode().IsRegular() && strings.HasSuffix(strings.ToLower(file.Name()), suffix) {
			filespec := filepath.Join(inputDir, file.Name())

			data, err := ioutil.ReadFile(filespec)
			if err != nil {
				return errors.New(err.Error())
			}

			err = processor(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
