..
..  Copyright (c) 2019 AT&T Intellectual Property.
..  Copyright (c) 2019 Nokia.
..
..  Licensed under the Creative Commons Attribution 4.0 International
..  Public License (the "License"); you may not use this file except
..  in compliance with the License. You may obtain a copy of the License at
..
..    https://creativecommons.org/licenses/by/4.0/
..
..  Unless required by applicable law or agreed to in writing, documentation
..  distributed under the License is distributed on an "AS IS" BASIS,
..  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
..
..  See the License for the specific language governing permissions and
..  limitations under the License.
..

Release-Notes
=============

This document provides the release notes of the dbaas.

.. contents::
   :depth: 3
   :local:



Version history
---------------

[0.4.1] - 2020-06-17

* Upgrade base image to bldr-alpine3:12-a3.11 in Redis docker build

[0.4.0] - 2020-04-23

* Bump version to 0.4.0 to follow RIC versioning rules (4 is meaning RIC release R4). No functional changes.

[0.3.2] - 2020-04-22

* Upgrade base image to bldr-alpine3:10-a3.22-rmr3 in Redis docker build
* Fix redismodule resource leak

[0.3.1] - 2020-02-13

* Upgrade base image to alpine3-go:1-rmr1.13.1 in Redis docker build

[0.3.0] - 2020-01-23

* Enable unit tests and valgrind in CI.
* Update redismodule with new commands.
* Update documentation.

[0.2.2] - 2019-11-12

* Take Alpine (version 6-a3.9) linux base image into use in Redis docker image.
* Add mandatory documentation files.

[0.2.1] - 2019-09-17

* Add the curl tool to docker image to facilitate trouble-shooting.

[0.2.0] - 2019-09-03

* Take Redis 5.0 in use.

[0.1.0] - 2019-06-17

* Initial Implementation to provide all the needed elements to deploy database
  as a service docker image to kubernetes.
* Introduce new Redis modules: SETIE, SETNE, DELIE, DELNE, MSETPUB, MSETMPUB,
  SETXXPUB, SETNXPUB, SETIEPUB, SETNEPUB, DELPUB, DELMPUB, DELIEPUB, DELNEPUB,
  NGET, NDEL.
