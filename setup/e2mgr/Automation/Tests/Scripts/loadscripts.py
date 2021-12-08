##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#


import redis
import config

def getRedisClient():

    c = config.redis_ip_address

    p = config.redis_ip_port

    return redis.Redis(host=c, port=p, db=0)


def verify_value():

    r = getRedisClient()

    value = "\b\x98\xf7\xdd\xa3\xc7\xb4\x83\xde\x15\x12\x11\n\x0f02f829:0007ab00"

    if r.get("{e2Manager},LOAD:test1") != value:
        return True
    else:
        return False


def add():

    r = getRedisClient()

    r.set("{e2Manager},LOAD:test1", "\b\x98\xf7\xdd\xa3\xc7\xb4\x83\xde\x15\x12\x11\n\x0f02f829:0007ab00")

    if r.exists("{e2Manager},LOAD:test1"):
        return True
    else:
        return False


def verify():

    r = getRedisClient()

    if r.exists("{e2Manager},LOAD:test1"):
        return True
    else:
        return False

