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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     Collections
Library     REST      ${url}
Resource    ../Resource/scripts_variables.robot
Library     String
Library     Process
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/rsmscripts.py



*** Test Cases ***

Prepare Ran in Connected connectionStatus
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED

Run Reset from RAN
    Run    ${Run_Config}
    Sleep   1s

Prepare logs for tests
    Remove log files
    Save logs

Verify logs - Reset Sent by simulator
    ${Reset}=   Grep File  ./${gnb_log_filename}  ResetRequest has been sent
    Should Be Equal     ${Reset}     gnbe2_simu: ResetRequest has been sent

Verify logs - e2mgr logs - messege sent
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_REQ_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

Verify logs - e2mgr logs - messege received
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_RESP_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

RAN Restarted messege sent
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RAN_RESTARTED_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

RSM RESOURCE STATUS REQUEST message sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test1}
    Should Be Equal As Strings    ${result}      True

Verify RSM RAN info exists in redis
   ${result}=   rsmscripts.verify_rsm_ran_info_start_false
   Should Be Equal As Strings  ${result}    True