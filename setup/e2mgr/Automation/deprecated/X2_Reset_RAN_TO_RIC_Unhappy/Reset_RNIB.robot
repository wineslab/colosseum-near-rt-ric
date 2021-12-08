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
Library     ../Scripts/find_error_script.py
Suite Teardown  Start Dbass with 4 dockers


*** Test Cases ***

Prepare Ran in Connected connectionStatus
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED

Stop RNIB
    Stop Dbass


Run Reset from RAN
    Run    ${Run_Config}
    Sleep   60s

Prepare logs for tests
    Remove log files
    Save logs

Verify logs - Reset Sent by simulator
    ${Reset}=   Grep File  ./${gnb_log_filename}  ResetRequest has been sent
    Should Be Equal     ${Reset}     gnbe2_simu: ResetRequest has been sent

Verify logs for restart received
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_REQ_message_type}    ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

Verify for error on retrying
    ${result}    find_error_script.find_error    ${EXECDIR}     ${e2mgr_log_filename}   ${failed_to_retrieve_nodeb_message}
    Should Be Equal As Strings    ${result}      True


*** Keywords ***
Start Dbass with 4 dockers
     Run And Return Rc And Output    ${dbass_remove}
     Run And Return Rc And Output    ${dbass_start}
     Sleep  5s
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number-1}

