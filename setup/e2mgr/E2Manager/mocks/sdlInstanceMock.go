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

import "github.com/stretchr/testify/mock"

type MockSdlInstance struct {
	mock.Mock
}

func (m *MockSdlInstance) SubscribeChannel(cb func(string, ...string), channels ...string) error {
	a := m.Called(cb, channels)
	return a.Error(0)
}

func (m *MockSdlInstance) UnsubscribeChannel(channels ...string) error {
	a := m.Called(channels)
	return a.Error(0)
}

func (m *MockSdlInstance) SetAndPublish(channelsAndEvents []string, pairs ...interface{}) error {
	a := m.Called(channelsAndEvents, pairs)
	return a.Error(0)
}

func (m *MockSdlInstance) SetIfAndPublish(channelsAndEvents []string, key string, oldData, newData interface{}) (bool, error) {
	a := m.Called(channelsAndEvents, key, oldData, newData)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlInstance) SetIfNotExistsAndPublish(channelsAndEvents []string, key string, data interface{}) (bool, error) {
	a := m.Called(channelsAndEvents, key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlInstance) RemoveAndPublish(channelsAndEvents []string, keys []string) error {
	a := m.Called(channelsAndEvents, keys)
	return a.Error(0)
}

func (m *MockSdlInstance) RemoveIfAndPublish(channelsAndEvents []string, key string, data interface{}) (bool, error) {
	a := m.Called(channelsAndEvents, key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlInstance) RemoveAllAndPublish(channelsAndEvents []string) error {
	a := m.Called(channelsAndEvents)
	return a.Error(0)
}

func (m *MockSdlInstance) Set(pairs ...interface{}) error {
	a := m.Called(pairs)
	return a.Error(0)
}

func (m *MockSdlInstance) Get(keys []string) (map[string]interface{}, error) {
	a := m.Called(keys)
	return a.Get(0).(map[string]interface{}), a.Error(1)
}

func (m *MockSdlInstance) GetAll() ([]string, error) {
	a := m.Called()
	return a.Get(0).([]string), a.Error(1)
}

func (m *MockSdlInstance) Close() error {
	a := m.Called()
	return a.Error(0)
}

func (m *MockSdlInstance) Remove(keys []string) error {
	a := m.Called(keys)
	return a.Error(0)
}

func (m *MockSdlInstance) RemoveAll() error {
	a := m.Called()
	return a.Error(0)
}

func (m *MockSdlInstance) SetIf(key string, oldData, newData interface{}) (bool, error) {
	a := m.Called(key, oldData, newData)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlInstance) SetIfNotExists(key string, data interface{}) (bool, error) {
	a := m.Called(key, data)
	return a.Bool(0), a.Error(1)
}
func (m *MockSdlInstance) RemoveIf(key string, data interface{}) (bool, error) {
	a := m.Called(key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlInstance) AddMember(group string, member ...interface{}) error{
	a := m.Called(group, member)
	return a.Error(0)
}

func (m *MockSdlInstance) RemoveMember(group string, member ...interface{}) error {
	a := m.Called(group, member)
	return a.Error(0)
}
func (m *MockSdlInstance) RemoveGroup(group string) error {
	a := m.Called(group)
	return a.Error(0)
}
func (m *MockSdlInstance) GetMembers(group string) ([]string, error) {
	a := m.Called(group)
	return a.Get(0).([]string), a.Error(1)
}
func (m *MockSdlInstance) IsMember(group string, member interface{}) (bool, error){
	a := m.Called(group, member)
	return a.Bool(0), a.Error(1)
}
func (m *MockSdlInstance) GroupSize(group string) (int64, error){
	a := m.Called(group,)
	return int64(a.Int(0)), a.Error(1)
}
