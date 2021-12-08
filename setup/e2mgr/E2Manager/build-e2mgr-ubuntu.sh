#!/bin/bash
##############################################################################
#
#   Copyright (c) 2020 AT&T Intellectual Property.
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

# Installs libraries and builds E2 manager
# Prerequisites:
#   Debian distro; e.g., Ubuntu
#   NNG shared library
#   golang (go); tested with version 1.12
#   current working directory is E2Manager
#   running with sudo privs, which is default in Docker

# Stop at first error and be verbose
set -eux

echo "--> e2mgr-build-ubuntu.sh"

# Install RMR from deb packages at packagecloud.io
rmr=rmr_4.0.2_amd64.deb
wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/$rmr/download.deb
#sudo
dpkg -i $rmr
rm $rmr
rmrdev=rmr-dev_4.0.2_amd64.deb
wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/$rmrdev/download.deb
#sudo
dpkg -i $rmrdev
rm $rmrdev

# required to find nng and rmr libs
export LD_LIBRARY_PATH=/usr/local/lib

# go installs tools like go-acc to $HOME/go/bin
# ubuntu minion path lacks go
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin

# install the go coverage tool helper
go get -v github.com/ory/go-acc

cd 3rdparty/asn1codec \
    && make clobber \
    && make \
    && cd ../..

go build -v app/main.go

# Execute UT and measure coverage
# cgocheck=2 enables expensive checks that should not miss any errors,
# but will cause your program to run slower.
# clobberfree=1 causes the garbage collector to clobber the memory content
# of an object with bad content when it frees the object.
# gcstoptheworld=1 disables concurrent garbage collection, making every
# garbage collection a stop-the-world event.
# Setting gcstoptheworld=2 also disables concurrent sweeping after the
# garbage collection finishes.
# Setting allocfreetrace=1 causes every allocation to be profiled and a
# stack trace printed on each object's allocation and free.
export GODEBUG=cgocheck=2,clobberfree=1,gcstoptheworld=2,allocfreetrace=0
# Static route table is provided in git repo
export RMR_SEED_RT=$(pwd)/router_test.txt
# writes to coverage.txt by default
# SonarCloud accepts the text format
go-acc $(go list ./... | grep -vE '(/mocks|/tests|/e2managererrors|/enums)' )

# TODO: drop rewrite of path prefix when SonarScanner is extended
# rewrite the module name to a directory name in the coverage report
# https://jira.sonarsource.com/browse/SONARSLANG-450
sed -i -e 's/^e2mgr/E2Manager/' coverage.txt

echo "--> e2mgr-build-ubuntu.sh ends"
