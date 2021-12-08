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
Reset - 400 http - 401 Corrupted json
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    String   response body connectionStatus    CONNECTED
    Set Headers     ${header}
    PUT    /v1/nodeb/test1/reset  {abc}
    Integer    response status   400
    Integer    response body errorCode  401
    String     response body errorMessage  corrupted json

Reset - 400 http - 401 Validation error
    Set Headers     ${header}
    PUT    /v1/nodeb/test1/reset    ${resetbadcausejson}
    Integer    response status   400
    Integer    response body errorCode  402
    String     response body errorMessage  Validation error





