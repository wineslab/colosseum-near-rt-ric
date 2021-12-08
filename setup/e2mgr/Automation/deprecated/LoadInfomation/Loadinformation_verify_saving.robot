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
Suite Setup  Prepare Simulator For Load Information
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     Process
Library     ../Scripts/loadscripts.py
Library     REST      ${url}
Suite Teardown  Stop Simulator





*** Test Cases ***
Verify Load information doesn't exist in redis
    ${result}=     loadscripts.verify
    Should Be Equal As Strings      ${result}   False


Trigger X-2 Setup for load information
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body connectionStatus    CONNECTED


Verify Load information does exist in redis
    Sleep  2s
    ${result}=     loadscripts.verify
    Should Be Equal As Strings      ${result}   True

