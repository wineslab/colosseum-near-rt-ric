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
Resource    ../Resource/scripts_variables.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/rsmscripts.py
Library     ../Scripts/e2mdbscripts.py
Library     REST      ${url}

*** Test Cases ***
X2 - Setup Test
    Post Request setup node b x-2
    Integer     response status       204

X2 - Get Nodeb
    Get Request Node B Enb test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body ip    ${ip_gnb_simu} 
    Integer  response body port     5577
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     ENB
    String   response body associatedE2tInstanceAddress     e2t.att.com:38000  
    String   response body enb enbType     MACRO_ENB
    Integer  response body enb servedCells 0 pci  99
    String   response body enb servedCells 0 cellId   02f829:0007ab00
    String   response body enb servedCells 0 tac    0102
    String   response body enb servedCells 0 broadcastPlmns 0   "02f829"
    Integer  response body enb servedCells 0 choiceEutraMode fdd ulearFcn    1
    Integer  response body enb servedCells 0 choiceEutraMode fdd dlearFcn    1
    String   response body enb servedCells 0 choiceEutraMode fdd ulTransmissionBandwidth   BW50
    String   response body enb servedCells 0 choiceEutraMode fdd dlTransmissionBandwidth   BW50


prepare logs for tests
    Remove log files
    Save logs


X2 - RAN Connected message going to be sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}   ${e2mgr_log_filename}  ${RAN_CONNECTED_message_type}    ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

RSM RESOURCE STATUS REQUEST message sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test1}
    Should Be Equal As Strings    ${result}      True

Verify RSM RAN info exists in redis
   ${result}=   rsmscripts.verify_rsm_ran_info_start_false
   Should Be Equal As Strings  ${result}    True

Verify RAN is associated with E2T instance
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance     test1    e2t.att.com:38000
   Should Be True    ${result}



