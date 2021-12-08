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

Overview
========

The ric-plt/dbaas repo provides all the needed elements to deploy database as
a service (Dbaas) to kubernetes. Dbaas service is realized with a single
container running Redis database. The database is configured to be
non-persistent and non-redundant.

For the time being Dbaas only allowed usage is to provide database backend
service for Shared Data Layer (SDL).
