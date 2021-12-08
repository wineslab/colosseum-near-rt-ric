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

*** Settings ***
Suite Setup  Prepare Enviorment
Resource    ../Resource/scripts_variables.robot
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/rsmscripts.py
Library     REST      ${url}


*** Test Cases ***
Prepare Ran in Connected connectionStatus
    Post Request setup node b endc-setup
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test2
    Integer  response status  200
    String   response body ranName    test2
    String   response body connectionStatus    CONNECTED


Send Reset reqeust with cause
    Set Headers     ${header}
    PUT    /v1/nodeb/test2/reset    ${resetcausejson}
    Integer  response status  204


Prepare logs for tests
    Remove log files
    Save logs


RAN Restarted messege sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${e2mgr_log_filename}  ${RAN_RESTARTED_message_type}    ${Meid_test2}
    Should Be Equal As Strings    ${result}      True

RSM RESOURCE STATUS REQUEST message not sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test2}
    Should Be Equal As Strings    ${result}      False

Verify RSM RAN info doesn't exist in redis
   ${result}=   rsmscripts.verify_rsm_ran_info_start_false
   Should Be Equal As Strings  ${result}    False