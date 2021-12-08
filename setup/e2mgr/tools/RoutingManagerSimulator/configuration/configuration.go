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


package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	Http struct {
		Port int
	}

}

func ParseConfiguration() *Configuration{
	viper.SetConfigType("yaml")
	viper.SetConfigName("configuration")
	viper.AddConfigPath("RoutingManagerSimulator/resources/")
	viper.AddConfigPath("./resources/")  //For production
	viper.AddConfigPath("../resources/") //For test under Docker
	viper.AddConfigPath("../../resources/") //For test under Docker
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("#configuration.ParseConfiguration - failed to read configuration file: %s\n", err))
	}

	config := Configuration{}
	config.fillHttpConfig(viper.Sub("http"))
	return &config
}

func (c *Configuration)fillHttpConfig(httpConfig *viper.Viper) {
	if httpConfig == nil {
		panic(fmt.Sprintf("#configuration.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n"))
	}
	c.Http.Port = httpConfig.GetInt("port")
}
