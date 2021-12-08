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
Library     Collections
Resource   ../Resource/Keywords.robot
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}





*** Test Cases ***

Pre Condition for Connecting - no simu
    Run And Return Rc And Output    ${stop_simu}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number-1}

Reset - 400 http - 403 wrong state
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  10s
    GET      /v1/nodeb/test1
    String   response body connectionStatus    DISCONNECTED
    Set Headers     ${header}
    PUT    /v1/nodeb/test1/reset
    #Output
    Integer    response status   400
    Integer    response body errorCode  403
    String     response body errorMessage   ${403_reset_message}




