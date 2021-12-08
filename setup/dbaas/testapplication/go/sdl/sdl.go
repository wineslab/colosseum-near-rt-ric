//   Copyright (c) 2019 AT&T Intellectual Property.
//   Copyright (c) 2019 Nokia.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

//
//   This source code is part of the near-RT RIC (RAN Intelligent Controller)
//   platform project (RICP).
//

package sdl

import (
	"github.com/go-redis/redis"
	"os"
	"reflect"
)

type SdlInstance struct {
	nameSpace string
	nsPrefix  string
	client    *redis.Client
}

func Create(nameSpace string) *SdlInstance {
	hostname := os.Getenv("DBAAS_SERVICE_HOST")
	if hostname == "" {
		hostname = "localhost"
	}
	port := os.Getenv("DBAAS_SERVICE_PORT")
	if port == "" {
		port = "6379"
	}
	redisAddress := hostname + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s := SdlInstance{
		nameSpace: nameSpace,
		nsPrefix:  "{" + nameSpace + "},",
		client:    client,
	}

	return &s
}

func (s *SdlInstance) setNamespaceToKeys(pairs ...interface{}) []interface{} {
	var retVal []interface{}
	for i, v := range pairs {
		if i%2 == 0 {
			reflectType := reflect.TypeOf(v)
			switch reflectType.Kind() {
			case reflect.Slice:
				x := reflect.ValueOf(v)
				for i2 := 0; i2 < x.Len(); i2++ {
					if i2%2 == 0 {
						retVal = append(retVal, s.nsPrefix+x.Index(i2).Interface().(string))
					} else {
						retVal = append(retVal, x.Index(i2).Interface())
					}
				}
			case reflect.Array:
				x := reflect.ValueOf(v)
				for i2 := 0; i2 < x.Len(); i2++ {
					if i2%2 == 0 {
						retVal = append(retVal, s.nsPrefix+x.Index(i2).Interface().(string))
					} else {
						retVal = append(retVal, x.Index(i2).Interface())
					}
				}
			default:
				retVal = append(retVal, s.nsPrefix+v.(string))
			}
		} else {
			retVal = append(retVal, v)
		}
	}
	return retVal
}

func (s *SdlInstance) Set(pairs ...interface{}) error {
	keyAndData := s.setNamespaceToKeys(pairs...)
	err := s.client.MSet(keyAndData...).Err()
	return err
}

func (s *SdlInstance) Get(keys []string) (map[string]interface{}, error) {
	var keysWithNs []string
	for _, v := range keys {
		keysWithNs = append(keysWithNs, s.nsPrefix+v)
	}
	val, err := s.client.MGet(keysWithNs...).Result()
	m := make(map[string]interface{})
	if err != nil {
		return m, err
	}
	for i, v := range val {
		m[keys[i]] = v
	}
	return m, err
}

func (s *SdlInstance) SetIf(key string, oldData, newData interface{}) {
	panic("SetIf not implemented\n")
}

func (s *SdlInstance) SetIfiNotExists(key string, data interface{}) {
	panic("SetIfiNotExists not implemented\n")
}

func (s *SdlInstance) Remove(keys ...string) {
	panic("Remove not implemented\n")
}

func (s *SdlInstance) RemoveIf(key string, data interface{}) {
	panic("RemoveIf not implemented\n")
}

func (s *SdlInstance) GetAll() []string {
	panic("GetAll not implemented\n")
}

func (s *SdlInstance) RemoveAll() {
	panic("RemoveAll not implemented\n")
}

