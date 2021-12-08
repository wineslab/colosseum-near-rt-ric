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
Resource    ../Resource/scripts_variables.robot
Resource   ../Resource/Keywords.robot
Resource   ../Resource/resource.robot
Library    ../Scripts/find_error_script.py
Library     OperatingSystem
Library     REST      ${url}
Suite Teardown   Start Dbass

*** Test Cases ***
Get node b gnb - DB down - 500
    Stop Dbass
    GET      /v1/nodeb/test5
    Integer  response status            500
    Integer  response body errorCode            500
    String   response body errorMessage     RNIB error


Prepare logs for tests
    Remove log files
    Save logs

Verify e2mgr logs - First retry to retrieve from db
  ${result}    find_error_script.find_error     ${EXECDIR}  ${e2mgr_log_filename}    ${first_retry_to_retrieve_from_db}
   Should Be Equal As Strings    ${result}      True

Verify e2mgr logs - Third retry to retrieve from db
   ${result}    find_error_script.find_error     ${EXECDIR}  ${e2mgr_log_filename}   ${third_retry_to_retrieve_from_db}
   Should Be Equal As Strings    ${result}      True

