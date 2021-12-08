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
Library     REST        ${url}




*** Test Cases ***

Update Ran
    Sleep  2s
    Update Ran request
    Integer  response status  200
    String   response body ranName    ${ranname}
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     GNB
    String   response body gnb servedNrCells 0 servedNrCellInformation cellId   abcd
    String   response body gnb servedNrCells 0 nrNeighbourInfos 0 nrCgi  one
    String   response body gnb servedNrCells 0 servedNrCellInformation servedPlmns 0  whatever











