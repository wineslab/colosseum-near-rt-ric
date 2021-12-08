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
Documentation   Keywords file
Library     ../Scripts/cleanup_db.py
Resource   ../Resource/resource.robot
Library     OperatingSystem

*** Keywords ***
Get Request node b gnb
    Sleep    1s
    GET      ${getNodeb}

Update Ran request
    Sleep  1s
    PUT    ${update_gnb_url}   ${update_gnb_body}


Update Ran request not valid
    Sleep  1s
    PUT    ${update_gnb_url}   ${update_gnb_body_notvalid}

Remove log files
    Remove File  ${EXECDIR}/${gnb_log_filename}
    Remove File  ${EXECDIR}/${e2mgr_log_filename}
    Remove File  ${EXECDIR}/${e2t_log_filename}
    Remove File  ${EXECDIR}/${rm_sim_log_filename}

Save logs
    Sleep   1s
    Run     ${Save_sim_log}
    Run     ${Save_e2mgr_log}
    Run     ${Save_e2t_log}
    Run     ${Save_rm_sim_log}

Stop Simulator
    Run And Return Rc And Output    ${stop_simu}

Prepare Enviorment
     Log To Console  Starting preparations
     ${starting_timestamp}    Evaluate   datetime.datetime.now(datetime.timezone.utc).isoformat("T")   modules=datetime 
     ${e2t_log_filename}      Evaluate      "e2t.${SUITE NAME}.log".replace(" ","-")
     ${e2mgr_log_filename}    Evaluate      "e2mgr.${SUITE NAME}.log".replace(" ","-")
     ${gnb_log_filename}      Evaluate      "gnb.${SUITE NAME}.log".replace(" ","-")
     ${rm_sim_log_filename}   Evaluate      "rm_sim.${SUITE NAME}.log".replace(" ","-")
     ${Save_sim_log}          Evaluate   'docker logs --since ${starting_timestamp} gnbe2_oran_simu > ${gnb_log_filename}'
     ${Save_e2mgr_log}        Evaluate   'docker logs --since ${starting_timestamp} e2mgr > ${e2mgr_log_filename}'
     ${Save_e2t_log}          Evaluate   'docker logs --since ${starting_timestamp} e2 > ${e2t_log_filename}'
     ${Save_rm_sim_log}       Evaluate   'docker logs --since ${starting_timestamp} rm_sim > ${rm_sim_log_filename}'
     Set Suite Variable  ${e2t_log_filename}
     Set Suite Variable  ${e2mgr_log_filename}  
     Set Suite Variable  ${gnb_log_filename}   
     Set Suite Variable  ${rm_sim_log_filename}
     Set Suite Variable  ${Save_sim_log}
     Set Suite Variable  ${Save_e2mgr_log}
     Set Suite Variable  ${Save_e2t_log}
     Set Suite Variable  ${Save_rm_sim_log}

	 Log To Console  Ready to flush db
     ${flush}  cleanup_db.flush
     Should Be Equal As Strings  ${flush}  True
     Run And Return Rc And Output    ${stop_simu}
     Run And Return Rc And Output    ${docker_Remove}
     Run And Return Rc And Output    ${run_simu_regular}
     Sleep  3s
     Log To Console  Validating dockers are up
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}

Start E2
     Run And Return Rc And Output    ${start_e2}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}
     Sleep  2s

Stop E2
     Run And Return Rc And Output    ${stop_e2}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number-1}
     Sleep  2s

Start Dbass
     Run And Return Rc And Output    ${dbass_remove}
     Run And Return Rc And Output    ${dbass_start}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}

Stop Dbass
     Run And Return Rc And Output    ${dbass_stop}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number-1}

Restart simulator
    Run And Return Rc And Output    ${restart_simu}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number}

Start RoutingManager Simulator
    Run And Return Rc And Output    ${start_routingmanager_sim}

Stop RoutingManager Simulator
    Run And Return Rc And Output    ${stop_routingmanager_sim}

Restart simulator with less docker
    Run And Return Rc And Output    ${restart_simu}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number-1}

