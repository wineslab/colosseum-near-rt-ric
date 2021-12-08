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
	"k8s.io/client-go/kubernetes"
)

type KubernetesManager struct {
	Logger    *logger.Logger
	ClientSet kubernetes.Interface
	Config    *configuration.Configuration
}

func NewKubernetesManager(logger *logger.Logger, config *configuration.Configuration) *KubernetesManager {
	return &KubernetesManager{
		Logger:    logger,
		//ClientSet: createClientSet(logger, config),
		Config:    config,
	}
}

/*func createClientSet(logger *logger.Logger, config *configuration.Configuration) kubernetes.Interface {

	absConfigPath,err := filepath.Abs(config.Kubernetes.ConfigPath)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}

	kubernetesConfig, err := clientcmd.BuildConfigFromFlags("", absConfigPath)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}

	clientSet, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		logger.Errorf("#KubernetesManager.init - error: %s", err)
		return nil
	}
	return clientSet
}*/

/*func (km KubernetesManager) DeletePod(podName string) error {
	km.Logger.Infof("#KubernetesManager.DeletePod - POD name: %s ", podName)

	if km.ClientSet == nil {
		km.Logger.Errorf("#KubernetesManager.DeletePod - no kubernetesManager connection")
		return e2managererrors.NewInternalError()
	}

	if len(podName) == 0 {
		km.Logger.Warnf("#KubernetesManager.DeletePod - empty pod name")
		return e2managererrors.NewInternalError()
	}

	err := km.ClientSet.CoreV1().Pods(km.Config.Kubernetes.KubeNamespace).Delete(podName, &metaV1.DeleteOptions{})

	if err != nil {
		km.Logger.Errorf("#KubernetesManager.DeletePod - POD %s can't be deleted, error: %s", podName, err)
		return err
	}

	km.Logger.Infof("#KubernetesManager.DeletePod - POD %s was deleted", podName)
	return nil
}*/