.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0
.. Copyright (C) 2020 AT&T


Installation Guide
==================

.. contents::
   :depth: 3
   :local:

Abstract
--------

This document describes how to install the HelloWorld (HW) xAPP. 

Version history

+--------------------+--------------------+--------------------+--------------------+
| **Date**           | **Ver.**           | **Author**         | **Comment**        |
|                    |                    |                    |                    |
+--------------------+--------------------+--------------------+--------------------+
| 2020-05-19         |1.0.0               |Shraboni Jana       | Bronze Release     |
|                    |                    |                    |                    |
+--------------------+--------------------+--------------------+--------------------+


Introduction
------------

This document provides guidelines on how to install and configure the HW xAPP in various environments/operating modes.
The audience of this document is assumed to have good knowledge in RIC Platform.


Preface
-------
The HW xAPP can be run directly as a Linux binary, as a docker image, or in a pod in a Kubernetes environment.  The first
two can be used for testing/evaluation. The last option is how an xAPP is deployed in the RAN Intelligent Controller environment.
This document covers all three methods.  




Software Installation and Deployment
------------------------------------
The build process assumes a Linux environment with a gcc (>= 4.0)  compatible compiler and  has been tested on Ubuntu. For building docker images,
the Docker environment must be present in the system.


Build Process
~~~~~~~~~~~~~
The HW xAPP can be either tested as a Linux binary or as a docker image.
   1. **Linux binary**: 
      The HW xAPP may be compiled and invoked directly. Pre-requisite software packages that must be installed prior to compiling are documented in the Dockerfile in the repository. README file in the repository mentions the steps to be followed to make "hw-xapp-main" binary.   
   
   2. **Docker Image**: From the root of the repository, run   *docker --no-cache build -t <image-name> ./* .


Deployment
~~~~~~~~~~
**Invoking  xAPP docker container directly** (not in RIC Kubernetes env.):
        xAPP descriptor(config-file.json) is available say under directory /home/test-config,  the docker image can be invoked as *docker run --net host -it --rm -v "/home/test-config:/opt/ric/config" --name  "HW-xAPP" <image>*. 


Testing 
--------

Unit tests for various modules of the HW xAPP are under the *test/* repository. HW xAPP uses Google Test Framework for unit testing. Currently, the unit tests must be compiled and executed  in a Linux environment (*make* in test directory) and Docker Image(*docker build -f Dockerfile-Unit-Tests .*). All software packages required for compiling the HW xAPP must be installed (as listed in the Dockerfile-Unit-Tests) for linux binary. 

