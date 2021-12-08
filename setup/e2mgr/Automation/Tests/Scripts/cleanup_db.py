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
import config
import redis
import time


def flush():

    c = config.redis_ip_address

    p = config.redis_ip_port

    r = redis.Redis(host=c, port=p, db=0)

    r.flushall()

    r.set("{e2Manager},E2TAddresses", "[\"10.0.2.15:38000\"]")

    r.set("{e2Manager},E2TInstance:10.0.2.15:38000","{\"address\":\"10.0.2.15:38000\",\"associatedRanList\":[],\"keepAliveTimestamp\":" + str(int((time.time()+2) * 1000000000)) + ",\"state\":\"ACTIVE\",\"deletionTimeStamp\":0}")

    return True




