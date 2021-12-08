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
/*

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

type TestStruct struct {
	description string
	namespace   string
	expected    []string
	objs        []runtime.Object
}

func initKubernetesManagerTest(t *testing.T) *KubernetesManager {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{}
	config.Kubernetes.KubeNamespace = "oran"
	config.Kubernetes.ConfigPath = "somePath"

	kubernetesManager := NewKubernetesManager(logger, config)

	return kubernetesManager
}

func TestDelete_NoPodName(t *testing.T) {
	test := TestStruct{
		description: "2 namespace, 2 pods in oran",
		namespace:   "oran",
		objs:        []runtime.Object{pod("oran", "POD_Test_1"), pod("oran", "POD_Test_2"), pod("some-namespace", "POD_Test_1")},
	}

	kubernetesManager := initKubernetesManagerTest(t)

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := kubernetesManager.DeletePod("")
		assert.NotNil(t, err)
	})
}

func TestDelete_NoPods(t *testing.T) {
	test := TestStruct{
		description: "no pods",
		namespace:   "oran",
		expected:    nil,
		objs:        nil,
	}

	kubernetesManager := initKubernetesManagerTest(t)

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := kubernetesManager.DeletePod("POD_Test")
		assert.NotNil(t, err)
	})
}

func TestDelete_PodExists(t *testing.T) {
	test := TestStruct{
		description: "2 namespace, 2 pods in oran",
		namespace:   "oran",
		objs:        []runtime.Object{pod("oran", "POD_Test_1"), pod("oran", "POD_Test_2"), pod("some-namespace", "POD_Test_1")},
	}

	kubernetesManager := initKubernetesManagerTest(t)

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := kubernetesManager.DeletePod("POD_Test_1")
		assert.Nil(t, err)
	})
}

func TestDelete_NoPodInNamespace(t *testing.T) {
	test := TestStruct{
		description: "2 namespace, 2 pods in oran",
		namespace:   "oran",
		objs:        []runtime.Object{pod("oran", "POD_Test_1"), pod("oran", "POD_Test_2"), pod("some-namespace", "POD_Test")},
	}

	kubernetesManager := initKubernetesManagerTest(t)

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := kubernetesManager.DeletePod("POD_Test")
		assert.NotNil(t, err)
	})
}

func TestDelete_NoNamespace(t *testing.T) {
	test := TestStruct{
		description: "No oran namespace",
		namespace:   "oran",
		objs:        []runtime.Object{pod("some-namespace", "POD_Test_1"), pod("some-namespace", "POD_Test_2"), pod("some-namespace", "POD_Test")},
	}

	kubernetesManager := initKubernetesManagerTest(t)

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := kubernetesManager.DeletePod("POD_Test")
		assert.NotNil(t, err)
	})
}

func pod(namespace, image string) *v1.Pod {

	return &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        image,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}
}
*/