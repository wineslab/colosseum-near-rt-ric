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

import config
import redis


def getRedisClientDecodeResponse():

    c = config.redis_ip_address

    p = config.redis_ip_port

    return redis.Redis(host=c, port=p, db=0, decode_responses=True)


def verify_e2t_addresses_key():

    r = getRedisClientDecodeResponse()
    
    value = "[\"10.0.2.15:38000\"]"

    return r.get("{e2Manager},E2TAddresses") == value


def verify_e2t_instance_key():

    r = getRedisClientDecodeResponse()

    e2_address = "\"address\":\"10.0.2.15:38000\""
    e2_associated_ran_list = "\"associatedRanList\":[]"
    e2_state = "\"state\":\"ACTIVE\""

    e2_db_instance = r.get("{e2Manager},E2TInstance:10.0.2.15:38000")

    if e2_db_instance.find(e2_address) < 0:
        return False
    if e2_db_instance.find(e2_associated_ran_list) < 0:
        return False
    if e2_db_instance.find(e2_state) < 0:
        return False

    return True