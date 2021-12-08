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
Suite Setup   Prepare Enviorment
Resource   ../Resource/Keywords.robot
Resource   ../Resource/resource.robot
Resource    ../Resource/scripts_variables.robot
Library     REST      ${url}
Library    RequestsLibrary
Library    Collections
Library    OperatingSystem
Library    json
Library     ../Scripts/e2mdbscripts.py


*** Test Cases ***

Get E2T instances
    ${result}    e2mdbscripts.populate_e2t_instances_in_e2m_db_for_get_e2t_instances_tc
    Create Session  getE2tInstances  ${url}
    ${headers}=  Create Dictionary    Accept=application/json
    ${resp}=    Get Request   getE2tInstances     /v1/e2t/list    headers=${headers}
    Should Be Equal As Strings   ${resp.status_code}    200
    Should Be Equal As Strings    ${resp.content}        [{"e2tAddress":"e2t.att.com:38000","ranNames":["test1","test2","test3"]},{"e2tAddress":"e2t.att.com:38001","ranNames":[]}]
    ${flush}  cleanup_db.flush












