*** Settings ***
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource    ../Resource/scripts_variables.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/cleanup_db.py
Library     ../Scripts/e2t_db_script.py

*** Test Cases ***

Test New E2T Send Init
    Stop E2
    ${result}=    cleanup_db.flush
    Should Be Equal As Strings  ${result}    True
    Start E2

prepare logs for tests
    Remove log files
    Save logs

E2M Logs - Verify RMR Message
    ${result}    find_rmr_message.verify_logs   ${EXECDIR}   ${e2mgr_log_filename}  ${E2_INIT_message_type}    ${None}
    Should Be Equal As Strings    ${result}      True

Verify E2T keys in DB
    ${result}=    e2t_db_script.verify_e2t_addresses_key
    Should Be Equal As Strings  ${result}    True

    ${result}=    e2t_db_script.verify_e2t_instance_key
    Should Be Equal As Strings  ${result}    True



