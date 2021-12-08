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
Documentation    Message types resource file


*** Variables ***

${E2_INIT_message_type}    MType: 1100
${Setup_failure_message_type}    MType: 12003
${first_retry_to_retrieve_from_db}      RnibDataService.retry - retrying 1 GetNodeb
${third_retry_to_retrieve_from_db}      RnibDataService.retry - after 3 attempts of GetNodeb
${RIC_RES_STATUS_REQ_message_type_successfully_sent}     Message type: 10090 - Successfully sent RMR message
${E2_TERM_KEEP_ALIVE_REQ_message_type_successfully_sent}     Message type: 1101 - Successfully sent RMR message


