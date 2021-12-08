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
CFLAGS=-Wall -Wpedantic -std=c11 -Og  -I./inc -I./src -I./e2ap_engine -DASN_DISABLE_OER_SUPPORT  -DASN_PDU_COLLECTION -D_POSIX_C_SOURCE=200809L -ggdb
export CFLAGS
OBJDIR=lib
LIB=$(OBJDIR)/libe2ap_codec.a
LIBSRC=configuration_update_wrapper.c x2setup_request_wrapper.c x2reset_request_wrapper.c x2reset_response_wrapper.c asn1codec_utils.c
LIBOBJ=$(addprefix $(OBJDIR)/,$(LIBSRC:.c=.o))
TESTX2SETUPREQUEST=tests/x2setup_request_wrapper_test
TESTCONFUPDATE=tests/configuration_update_wrapper_test
TESTX2RESETREQUEST=tests/x2reset_request_wrapper_test
TESTX2RESETRESPONSE=tests/x2reset_response_wrapper_test
TESTUNPACKXER=tests/unpack_xer

.PHONY: all clean e2ap_engine

all: $(LIB) $(TESTX2SETUPREQUEST) $(TESTCONFUPDATE)  $(TESTX2RESETREQUEST) $(TESTUNPACKXER) $(TESTX2RESETRESPONSE)

e2ap_engine/libasncodec.a:
	cd e2ap_engine/ && make -f converter-example.mk

$(OBJDIR)/%.o: src/%.c e2ap_engine/*.h
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) -c $< -o $@

$(LIB): e2ap_engine/libasncodec.a $(LIBOBJ) 
	$(AR) -crv $(LIB) $(LIBOBJ) 

$(TESTX2SETUPREQUEST): % : $(LIB) src/%.c 
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) src/$@.c -o $@ $(LIB)  e2ap_engine/libasncodec.a

$(TESTCONFUPDATE): % : $(LIB) src/%.c 
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) src/$@.c -o $@ $(LIB)  e2ap_engine/libasncodec.a

$(TESTX2RESETREQUEST): % : $(LIB) src/%.c
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) src/$@.c -o $@ $(LIB)  e2ap_engine/libasncodec.a

$(TESTUNPACKXER): % : $(LIB) src/%.c
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) src/$@.c -o $@ $(LIB)  e2ap_engine/libasncodec.a

$(TESTX2RESETRESPONSE): % : $(LIB) src/%.c
	mkdir -p $(dir $@)
	$(CC) $(CFLAGS) src/$@.c -o $@ $(LIB)  e2ap_engine/libasncodec.a

clean:
	rm -rf $(OBJDIR) tests 
clobber:
	cd e2ap_engine/ && make -f converter-example.mk clean
	rm -rf $(OBJDIR) tests 
