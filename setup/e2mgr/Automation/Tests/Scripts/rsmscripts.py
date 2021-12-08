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
import json


def getRedisClientDecodeResponse():

    c = config.redis_ip_address

    p = config.redis_ip_port

    return redis.Redis(host=c, port=p, db=0, decode_responses=True)

def set_general_config_resource_status_false():

    r = getRedisClientDecodeResponse()
    r.set("{rsm},CFG:GENERAL:v1.0.0" , "{\"enableResourceStatus\":false,\"partialSuccessAllowed\":true,\"prbPeriodic\":true,\"tnlLoadIndPeriodic\":true,\"wwLoadIndPeriodic\":true,\"absStatusPeriodic\":true,\"rsrpMeasurementPeriodic\":true,\"csiPeriodic\":true,\"periodicityMs\":1,\"periodicityRsrpMeasurementMs\":3,\"periodicityCsiMs\":3}")

def verify_rsm_ran_info_start_false():

    r = getRedisClientDecodeResponse()
    
    value = "{\"ranName\":\"test1\",\"enb1MeasurementId\":1,\"enb2MeasurementId\":0,\"action\":\"start\",\"actionStatus\":false}"

    return r.get("{rsm},RAN:test1") == value


def verify_rsm_ran_info_start_true():

    r = getRedisClientDecodeResponse()

    rsmInfoStr = r.get("{rsm},RAN:test1")
    rsmInfoJson = json.loads(rsmInfoStr)

    response = rsmInfoJson["ranName"] == "test1" and rsmInfoJson["enb1MeasurementId"] == 1 and rsmInfoJson["enb2MeasurementId"] != 1 and rsmInfoJson["action"] == "start" and rsmInfoJson["actionStatus"] == True

    return response


def verify_rsm_ran_info_stop_false():

    r = getRedisClientDecodeResponse()

    rsmInfoStr = r.get("{rsm},RAN:test1")
    rsmInfoJson = json.loads(rsmInfoStr)

    response = rsmInfoJson["ranName"] == "test1" and rsmInfoJson["enb1MeasurementId"] == 1 and rsmInfoJson["action"] == "stop" and rsmInfoJson["actionStatus"] == False

    return response


def verify_rsm_ran_info_stop_true():

    r = getRedisClientDecodeResponse()

    rsmInfoStr = r.get("{rsm},RAN:test1")
    rsmInfoJson = json.loads(rsmInfoStr)

    response = rsmInfoJson["ranName"] == "test1" and rsmInfoJson["action"] == "stop" and rsmInfoJson["actionStatus"] == True

    return response

def verify_general_config_enable_resource_status_true():

    r = getRedisClientDecodeResponse()

    configStr = r.get("{rsm},CFG:GENERAL:v1.0.0")
    configJson = json.loads(configStr)

    return configJson["enableResourceStatus"] == True

def verify_general_config_enable_resource_status_false():

    r = getRedisClientDecodeResponse()

    configStr = r.get("{rsm},CFG:GENERAL:v1.0.0")
    configJson = json.loads(configStr)

    return configJson["enableResourceStatus"] == False