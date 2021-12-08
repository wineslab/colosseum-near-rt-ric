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
Resource   ../Resource/Keywords.robot
Resource   ../Resource/resource.robot
Library     ../Scripts/e2mdbscripts.py
Library     REST      ${url}
Suite Teardown   Start RoutingManager Simulator

*** Test Cases ***
ENDC-setup - 503 http - 511 No Routing Manager Available
    Stop RoutingManager Simulator
    Set Headers     ${header}
    POST     /v1/nodeb/x2-setup    ${json}
    Integer  response status            503
    Integer  response body errorCode            511
    String   response body errorMessage     No Routing Manager Available

Verify RAN is NOT associated with E2T instance
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance     test1    e2t.att.com:38000
   Should Be True    ${result} == False
