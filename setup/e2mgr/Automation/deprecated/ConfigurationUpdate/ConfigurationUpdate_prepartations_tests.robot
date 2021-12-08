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
Library     REST      ${url}
Resource    ../Resource/scripts_variables.robot
Library     String
Library     Process
Library     ../Scripts/find_rmr_message.py


*** Test Cases ***
X2 - Setup and Get
    Post Request setup node b x-2
    Get Request node b enb test1
    String   response body connectionStatus    CONNECTED


Run Configuration update
    Run    ${Run_Config}
    Sleep   1s

Prepare logs for tests
    Remove log files
    Save logs

Verify logs - Confiugration update - Begin Tag Get
    ${Configuration}=   Grep File  ./${gnb_log_filename}  <ENDCConfigurationUpdate>
    ${ConfigurationAfterStrip}=     Strip String    ${Configuration}
    Should Be Equal     ${ConfigurationAfterStrip}        <ENDCConfigurationUpdate>

Verify logs - Confiugration update - End Tag Get
    ${ConfigurationEnd}=   Grep File  ./${gnb_log_filename}  </ENDCConfigurationUpdate>
    ${ConfigurationEndAfterStrip}=     Strip String    ${ConfigurationEnd}
    Should Be Equal     ${ConfigurationEndAfterStrip}        </ENDCConfigurationUpdate>

Verify logs - Confiugration update - Ack Tag Begin
    ${ConfigurationAck}=   Grep File  ./${gnb_log_filename}   <ENDCConfigurationUpdateAcknowledge>
    ${ConfigurationAckAfter}=     Strip String    ${ConfigurationAck}
    Should Be Equal     ${ConfigurationAckAfter}        <ENDCConfigurationUpdateAcknowledge>

Verify logs - Confiugration update - Ack Tag End
    ${ConfigurationAckEnd}=   Grep File  ./${gnb_log_filename}  </ENDCConfigurationUpdateAcknowledge>
    ${ConfigurationAckEndAfterStrip}=     Strip String    ${ConfigurationAckEnd}
    Should Be Equal     ${ConfigurationAckEndAfterStrip}        </ENDCConfigurationUpdateAcknowledge>

Verify logs - find RIC_ENDC_CONF_UPDATE
   ${result}   find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${configurationupdate_message_type}  ${Meid_test1}
   Should Be Equal As Strings    ${result}      True
Verify logs - find RIC_ENDC_CONF_UPDATE_ACK
   ${result1}  find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${configurationupdate_ack_message_type}  ${Meid_test1}
   Should Be Equal As Strings    ${result1}      True






