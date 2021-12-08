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

import config
import redis
import cleanup_db


def add():

    c = config.redis_ip_address

    p = config.redis_ip_port

    r = redis.Redis(host=c, port=p, db=0)

    cleanup_db.flush()

    r.set("{e2Manager},ENB:02f829:007a80", "\n\x05test1\x12\t10.0.2.15\x18\xc9+ \x01*\x10\n\x0602f829\x12\x06007a800\x01:3\b\x01\x12/\bc\x12\x0f02f829:0007ab50\x1a\x040102\"\x0602f829*\n\n\b\b\x01\x10\x01\x18\x04 \x040\x01")

    r.set("{e2Manager},RAN:test1","\x12\t10.0.2.15\x18\xc9+ \x03H\x01R\x02\b\t")

    r.set("{e2Manager},PCI:test1:63" , "\b\x01\x12/\bc\x12\x0f02f829:0007ab50\x1a\x040102\"\x0602f829*\n\n\b\b\x01\x10\x01\x18\x04 \x040\x01")

    r.set("{e2Manager},CELL:02f829:0007ab50" ,  "\b\x01\x12/\bc\x12\x0f02f829:0007ab50\x1a\x040102\"\x0602f829*\n\n\b\b\x01\x10\x01\x18\x04 \x040\x01")

    r.sadd("{e2Manager},ENB" , "\n\x05test1\x12\x10\n\x0602f829\x12\x06007a80")


    r.set("{e2Manager},GNB:03f829:002234", "\n\x05test2\x12\t10.0.2.16\x18\xc9+ \x01*\x10\n\x0702f829\x12\x070012340\x02BI\nG\nE\bc\x12\x1102f829:0008ab0120*\x0602f8290\x01:$\n\"\n\t\bd\"\x05\b\t\x12\x01\t\x12\t\bd\"\x05\b\t\x12\x01\t\x1a\x04\b\x01\x10\x01\"\x04\b\x01\x10\x01")

    r.set("{e2Manager},RAN:test2", "\n\x05test2\x12\t10.0.2.15\x18\xc9+ \x01*\x10\n\x0702f829\x12\x070012340\x03BI\nG\nE\bc\x12\x1103f829:0008ab0120*\x0602f8290\x01:$\n\"\n\t\bd\"\x05\b\t\x12\x01\t\x12\t\bd\"\x05\b\t\x12\x01\t\x1a\x04\b\x01\x10\x01\"\x04\b\x01\x10\x01")

    r.set("{e2Manager},PCI:test2:63", "\b\x02\x1aG\nE\bc\x12\x1102f829:0008ab0120*\x0702f8290\x01:$\n\"\n\t\bd\"\x05\b\t\x12\x01\t\x12\t\bd\"\x05\b\t\x12\x01\t\x1a\x04\b\x01\x10\x01\"\x04\b\x01\x10\x01")

    r.set("{e2Manager},NRCELL:02f829:0007ab0120", "\b\x02\x1aG\nE\bc\x12\x1102f829:0007ab0120*\x0602f8290\x01:$\n\"\n\t\bd\"\x05\b\t\x12\x01\t\x12\t\bd\"\x05\b\t\x12\x01\t\x1a\x04\b\x01\x10\x01\"\x04\b\x01\x10\x01")

    r.sadd("{e2Manager},GNB","\n\x05test2\x12\x10\n\x0603f829\x12\x06001234")

    return True
