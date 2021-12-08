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
Library     REST      ${url}
Suite Teardown   Start Dbass

*** Test Cases ***
Get All nodes - 500 http - 500 RNIB error
    Stop Dbass
    GET      /v1/nodeb/ids
    Integer  response status            500
    Integer  response body errorCode            500
    String   response body errorMessage     RNIB error


