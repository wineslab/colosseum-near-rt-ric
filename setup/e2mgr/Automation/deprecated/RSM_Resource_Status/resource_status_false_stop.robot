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
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource   ../Resource/scripts_variables.robot
Resource   resource_status_keywords.robot
Library    ../Scripts/rsmscripts.py
Library     ../Scripts/find_rmr_message.py
Library    OperatingSystem
Library    REST      ${url_rsm}
Suite Teardown  Delete All Sessions


*** Test Cases ***
Run setup
    rsmscripts.set_general_config_resource_status_false

    Prepare Ran In Connected Status

Put Http Stop Request To RSM
    Put Request Resource Status Stop
    Integer  response status  204

Verify RSM RAN Info Status Is Stop And True In Redis
    ${result}=   rsmscripts.verify_rsm_ran_info_stop_true
    Should Be Equal As Strings  ${result}    True

Verify RSM Enable Resource Status Is False In General Configuration In Redis
    ${result}=   rsmscripts.verify_general_config_enable_resource_status_false
    Should Be Equal As Strings  ${result}    True
    
prepare logs for tests
    Remove log files
    Save logs

Verify RSM Resource Status Request Message Not Sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test1}
    Should Be Equal As Strings    ${result}      False