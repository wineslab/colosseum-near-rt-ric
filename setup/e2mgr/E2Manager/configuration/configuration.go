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
	Logging struct {
		LogLevel string
	}
	Http struct {
		Port int
	}
	Rmr struct {
		Port       int
		MaxMsgSize int
	}
	RoutingManager struct {
		BaseUrl string
	}
/*	Kubernetes struct {
		ConfigPath string
		KubeNamespace  string
	}*/
	NotificationResponseBuffer   int
	BigRedButtonTimeoutSec       int
	MaxRnibConnectionAttempts    int
	RnibRetryIntervalMs          int
	KeepAliveResponseTimeoutMs   int
	KeepAliveDelayMs             int
	E2TInstanceDeletionTimeoutMs int
	GlobalRicId                  struct {
		PlmnId      string
		RicNearRtId string
	}
}

func ParseConfiguration() *Configuration {
	viper.SetConfigType("yaml")
	viper.SetConfigName("configuration")
	viper.AddConfigPath("E2Manager/resources/")
	viper.AddConfigPath("./resources/")     //For production
	viper.AddConfigPath("../resources/")    //For test under Docker
	viper.AddConfigPath("../../resources/") //For test under Docker
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("#configuration.ParseConfiguration - failed to read configuration file: %s\n", err))
	}

	config := Configuration{}
	config.populateRmrConfig(viper.Sub("rmr"))
	config.populateHttpConfig(viper.Sub("http"))
	config.populateLoggingConfig(viper.Sub("logging"))
	config.populateRoutingManagerConfig(viper.Sub("routingManager"))
	//config.populateKubernetesConfig(viper.Sub("kubernetes"))
	config.NotificationResponseBuffer = viper.GetInt("notificationResponseBuffer")
	config.BigRedButtonTimeoutSec = viper.GetInt("bigRedButtonTimeoutSec")
	config.MaxRnibConnectionAttempts = viper.GetInt("maxRnibConnectionAttempts")
	config.RnibRetryIntervalMs = viper.GetInt("rnibRetryIntervalMs")
	config.KeepAliveResponseTimeoutMs = viper.GetInt("keepAliveResponseTimeoutMs")
	config.KeepAliveDelayMs = viper.GetInt("KeepAliveDelayMs")
	config.E2TInstanceDeletionTimeoutMs = viper.GetInt("e2tInstanceDeletionTimeoutMs")
	config.populateGlobalRicIdConfig(viper.Sub("globalRicId"))
	return &config
}

func (c *Configuration) populateLoggingConfig(logConfig *viper.Viper) {
	if logConfig == nil {
		panic(fmt.Sprintf("#configuration.populateLoggingConfig - failed to populate logging configuration: The entry 'logging' not found\n"))
	}
	c.Logging.LogLevel = logConfig.GetString("logLevel")
}

func (c *Configuration) populateHttpConfig(httpConfig *viper.Viper) {
	if httpConfig == nil {
		panic(fmt.Sprintf("#configuration.populateHttpConfig - failed to populate HTTP configuration: The entry 'http' not found\n"))
	}
	c.Http.Port = httpConfig.GetInt("port")
}

func (c *Configuration) populateRmrConfig(rmrConfig *viper.Viper) {
	if rmrConfig == nil {
		panic(fmt.Sprintf("#configuration.populateRmrConfig - failed to populate RMR configuration: The entry 'rmr' not found\n"))
	}
	c.Rmr.Port = rmrConfig.GetInt("port")
	c.Rmr.MaxMsgSize = rmrConfig.GetInt("maxMsgSize")
}

func (c *Configuration) populateRoutingManagerConfig(rmConfig *viper.Viper) {
	if rmConfig == nil {
		panic(fmt.Sprintf("#configuration.populateRoutingManagerConfig - failed to populate Routing Manager configuration: The entry 'routingManager' not found\n"))
	}
	c.RoutingManager.BaseUrl = rmConfig.GetString("baseUrl")
}

/*func (c *Configuration) populateKubernetesConfig(rmConfig *viper.Viper) {
	if rmConfig == nil {
		panic(fmt.Sprintf("#configuration.populateKubernetesConfig - failed to populate Kubernetes configuration: The entry 'kubernetes' not found\n"))
	}
	c.Kubernetes.ConfigPath = rmConfig.GetString("configPath")
	c.Kubernetes.KubeNamespace = rmConfig.GetString("kubeNamespace")
}*/

func (c *Configuration) populateGlobalRicIdConfig(globalRicIdConfig *viper.Viper) {
	if globalRicIdConfig == nil {
		panic(fmt.Sprintf("#configuration.populateGlobalRicIdConfig - failed to populate Global RicId configuration: The entry 'globalRicId' not found\n"))
	}
	c.GlobalRicId.PlmnId = globalRicIdConfig.GetString("plmnId")
	c.GlobalRicId.RicNearRtId = globalRicIdConfig.GetString("ricNearRtId")
}

func (c *Configuration) String() string {
	return fmt.Sprintf("{logging.logLevel: %s, http.port: %d, rmr: { port: %d, maxMsgSize: %d}, routingManager.baseUrl: %s, "+
		"notificationResponseBuffer: %d, bigRedButtonTimeoutSec: %d, maxRnibConnectionAttempts: %d, "+
		"rnibRetryIntervalMs: %d, keepAliveResponseTimeoutMs: %d, keepAliveDelayMs: %d, e2tInstanceDeletionTimeoutMs: %d, "+
		"globalRicId: { plmnId: %s, ricNearRtId: %s}",//, kubernetes: {configPath: %s, kubeNamespace: %s}}",
		c.Logging.LogLevel,
		c.Http.Port,
		c.Rmr.Port,
		c.Rmr.MaxMsgSize,
		c.RoutingManager.BaseUrl,
		c.NotificationResponseBuffer,
		c.BigRedButtonTimeoutSec,
		c.MaxRnibConnectionAttempts,
		c.RnibRetryIntervalMs,
		c.KeepAliveResponseTimeoutMs,
		c.KeepAliveDelayMs,
		c.E2TInstanceDeletionTimeoutMs,
		c.GlobalRicId.PlmnId,
		c.GlobalRicId.RicNearRtId,
/*		c.Kubernetes.ConfigPath,
		c.Kubernetes.KubeNamespace,*/
	)
}
