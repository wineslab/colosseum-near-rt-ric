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

*** Settings ***
Documentation   Keywords file
Library     ../Scripts/cleanup_db.py
Resource   ../Resource/resource.robot
Library    Collections
Library    OperatingSystem
Library    json
Library    RequestsLibrary



*** Keywords ***
Prepare Ran In Connected Status
    Create Session  x2setup  ${url}
    ${headers}=  Create Dictionary    Accept=application/json    Content-Type=application/json
    ${resp}=    Post Request   x2setup    /v1/nodeb/x2-setup    data=${json_setup_rsm_tests}    headers=${headers}
    Should Be Equal As Strings    ${resp.status_code}    204

