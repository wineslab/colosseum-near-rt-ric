

#/*****************************************************************************
#                                                                            *
# Copyright 2020 AT&T Intellectual Property                                  *
#                                                                            *
# Licensed under the Apache License, Version 2.0 (the "License");            *
# you may not use this file except in compliance with the License.           *
# You may obtain a copy of the License at                                    *
#                                                                            *
#      http://www.apache.org/licenses/LICENSE-2.0                            *
#                                                                            *
# Unless required by applicable law or agreed to in writing, software        *
# distributed under the License is distributed on an "AS IS" BASIS,          *
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   *
# See the License for the specific language governing permissions and        *
# limitations under the License.                                             *
#                                                                            *
#******************************************************************************/

# For clarity: this generates object, not a lib as the CM command implies.
#
file(GLOB SOURCES "*.c")
file(GLOB HEADERS "*.h")

include_directories(.)
add_library( asn1_objects OBJECT ${SOURCES})

target_include_directories (asn1_objects PUBLIC
  $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
  $<INSTALL_INTERFACE:include>
  PRIVATE src)

if( DEV_PKG )                                   
  install( FILES
    ${HEADERS}
    DESTINATION ${install_inc}
    )
endif()
# ---------------------------------------------------------------------------
