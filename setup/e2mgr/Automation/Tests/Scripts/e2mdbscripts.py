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
import cleanup_db
import json

def getRedisClientDecodeResponse():
    c = config.redis_ip_address
    p = config.redis_ip_port
    return redis.Redis(host=c, port=p, db=0, decode_responses=True)

def verify_ran_is_associated_with_e2t_instance(ranName, e2tAddress):
    r = getRedisClientDecodeResponse()
    e2tInstanceJson = r.get("{e2Manager},E2TInstance:"+e2tAddress)
    e2tInstanceDic = json.loads(e2tInstanceJson)
    assocRanList = e2tInstanceDic.get("associatedRanList")
    return ranName in assocRanList

def verify_e2t_instance_has_no_associated_rans(e2tAddress):
    r = getRedisClientDecodeResponse()
    e2tInstanceJson = r.get("{e2Manager},E2TInstance:"+e2tAddress)
    e2tInstanceDic = json.loads(e2tInstanceJson)
    assocRanList = e2tInstanceDic.get("associatedRanList")
    return not assocRanList

def verify_e2t_instance_exists_in_addresses(e2tAddress):
    r = getRedisClientDecodeResponse()
    e2tAddressesJson = r.get("{e2Manager},E2TAddresses")
    e2tAddresses = json.loads(e2tAddressesJson)
    return e2tAddress in e2tAddresses

def verify_e2t_instance_key_exists(e2tAddress):
    r = getRedisClientDecodeResponse()
    return r.exists("{e2Manager},E2TInstance:"+e2tAddress)

def populate_e2t_instances_in_e2m_db_for_get_e2t_instances_tc():
    r = getRedisClientDecodeResponse()
    r.set("{e2Manager},E2TAddresses", "[\"e2t.att.com:38000\",\"e2t.att.com:38001\"]")
    r.set("{e2Manager},E2TInstance:e2t.att.com:38000", "{\"address\":\"e2t.att.com:38000\",\"associatedRanList\":[\"test1\",\"test2\",\"test3\"],\"keepAliveTimestamp\":1577619310484022369,\"state\":\"ACTIVE\"}")
    r.set("{e2Manager},E2TInstance:e2t.att.com:38001", "{\"address\":\"e2t.att.com:38001\",\"associatedRanList\":[],\"keepAliveTimestamp\":1577619310484022369,\"state\":\"ACTIVE\"}")
    return True

# def dissociate_ran_from_e2tInstance(ranName, e2tAddress):
#     r = getRedisClientDecodeResponse()
#     e2tInstanceJson = r.get("{e2Manager},E2TInstance:"+e2tAddress)
#     e2tInstanceDic = json.loads(e2tInstanceJson)
#     assocRanList = e2tInstanceDic.get("associatedRanList")
#     print(assocRanList)
#     assocRanList.remove(ranName)
#     updatedE2tInstanceJson = json.dumps(e2tInstanceDic)
#     print(updatedE2tInstanceJson)
#     r.set("{e2Manager},E2TInstance:"+e2tAddress, updatedE2tInstanceJson)
#     nodebBytes = r.get("{e2Manager},RAN:"+ranName)
#     encoded = nodebBytes.decode().replace(e2tAddress,"").encode()
#     r.set("{e2Manager},RAN:"+ranName, encoded)

