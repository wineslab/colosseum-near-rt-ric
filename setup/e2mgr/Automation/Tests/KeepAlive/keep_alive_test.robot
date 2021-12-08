robot##############################################################################
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
*** Settings ***
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource   ../Resource/scripts_variables.robot
Library    ../Scripts/find_error_script.py
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/rsmscripts.py
Library     ../Scripts/e2mdbscripts.py
Library    OperatingSystem
Library    Collections
Library     REST      ${url}


*** Test Cases ***

Get request gnb
    Sleep    2s
    Get Request node b gnb
    Integer  response status  200
    String   response body ranName    ${ranname}
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     GNB
    String   response body associatedE2tInstanceAddress  ${e2tinstanceaddress}
    Integer  response body gnb ranFunctions 0 ranFunctionId  1
    Integer  response body gnb ranFunctions 0 ranFunctionRevision  1
    Integer  response body gnb ranFunctions 1 ranFunctionId  2
    Integer  response body gnb ranFunctions 1 ranFunctionRevision  1
    Integer  response body gnb ranFunctions 2 ranFunctionId  3
    Integer  response body gnb ranFunctions 2 ranFunctionRevision  1


prepare logs for tests
    Remove log files
    Save logs

Verify RAN is associated with E2T instance
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance      ${ranname}    ${e2tinstanceaddress}
   Should Be True    ${result}

Stop E2T
    Stop E2
    Sleep  3s

Prepare logs
    Remove log files
    Save logs

Verify RAN is not associated with E2T instance
    Get Request node b gnb
    Integer  response status  200
    String   response body ranName    ${ranname}
    Missing  response body associatedE2tInstanceAddress
    String   response body connectionStatus    DISCONNECTED

Verify E2T instance removed from db
    ${result}    e2mdbscripts.verify_e2t_instance_key_exists     ${e2tinstanceaddress}
    Should Be True    ${result} == False

    ${result}    e2mdbscripts.verify_e2t_instance_exists_in_addresses     ${e2tinstanceaddress}
    Should Be True    ${result} == False

Start E2T
    Start E2