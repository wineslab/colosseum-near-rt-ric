CXX:= g++ --std=c++14 -O2 -L/usr/local/lib
CC:= gcc -O2 -L/usr/local/lib

TEST:=./
SRC:=../src
UTILSRC=../src/xapp-utils
MGMTSRC:=../src/xapp-mgmt

ASNSRC:=../asn1c_defs
E2APSRC:=../src/xapp-asn/e2ap
E2SMSRC:=../src/xapp-asn/e2sm

####### Logging library and flags
CLOGFLAGS:= `pkg-config mdclog --cflags`
LOG_LIBS:= `pkg-config mdclog --libs`
CURL_LIBS:= `pkg-config libcurl --libs`
RNIB_LIBS:= -pthread /usr/local/include/rnib/rnibreader.a
GTEST_LIBS:= -lgtest -lgtest_main
######## Keep include dirs separate so we have transparency


BASEFLAGS=  -Wall -std=c++14 $(CLOGFLAGS) 
C_BASEFLAGS= -Wall $(CLOGFLAGS) -DASN_DISABLE_OER_SUPPORT

XAPPFLAGS= -I$(SRC)
UTILFLAGS= -I$(UTILSRC)
MGMTFLAGS= -I$(MGMTSRC)

ASNFLAGS=-I$(ASNSRC) -DASN_DISABLE_OER_SUPPORT
E2APFLAGS = -I$(E2APSRC)
E2SMFLAGS = -I$(E2SMSRC)

########libs

LIBS= -lsdl -lrmr_si -lpthread -lm $(LOG_LIBS) $(CURL_LIBS) $(RNIB_LIBS) $(GTEST_LIBS)
COV_FLAGS= -fprofile-arcs -ftest-coverage

#######
TEST_SRC= hw_unit_tests.cc
XAPP_SRC= ../src/xapp.cc
TEST_HDR=$(wildcard $(TEST)/*.h)
UTIL_SRC= $(wildcard $(UTILSRC)/*.cc)
MGMT_SRC= $(wildcard $(MGMTSRC)/*.cc)

E2AP_SRC= $(wildcard $(E2APSRC)/*.cc)
E2SM_SRC= $(wildcard $(E2SMSRC)/*.cc)
ASN1C_SRC= $(wildcard $(ASNSRC)/*.c)

##############Objects
UTIL_OBJ=${UTIL_SRC:.cc=.o}
MGMT_OBJ=${MGMT_SRC:.cc=.o}
XAPP_OBJ=${XAPP_SRC:.cc=.o}
TEST_OBJ=${TEST_SRC:.cc=.o} 

E2AP_OBJ = $(E2AP_SRC:.cc=.o)
E2SM_OBJ = $(E2SM_SRC:.cc=.o)

ASN1C_MODULES = $(ASN1C_SRC:.c=.o)
$(ASN1C_MODULES): export CFLAGS = $(C_BASEFLAGS) $(ASNFLAGS)

$(UTIL_OBJ):export CPPFLAGS=$(BASEFLAGS) $(UTILFLAGS) $(MGMTFLAGS)
$(MGMT_OBJ):export CPPFLAGS=$(BASEFLAGS) $(MGMTFLAGS) $(ASNFLAGS) $(E2APFLAGS) $(E2SMFLAGS)

$(E2AP_OBJ): export CPPFLAGS = $(BASEFLAGS) $(ASNFLAGS) $(E2APFLAGS)
$(E2SM_OBJ): export CPPFLAGS = $(BASEFLAGS) $(ASNFLAGS) $(E2SMFLAGS)
$(XAPP_OBJ): export CPPFLAGS = $(BASEFLAGS) $(XAPPFLAGS) $(UTILFLAGS) $(MGMTFLAGS) $(E2APFLAGS) $(E2SMFLAGS) $(ASNFLAGS)

$(TEST_OBJ):export CPPFLAGS=$(BASEFLAGS) $(XAPPFLAGS) $(UTILFLAGS) $(MGMTFLAGS) $(E2APFLAGS) $(E2SMFLAGS) $(ASNFLAGS)
$(TEST_OBJ) = $(TEST_HDR) $(TEST_OBJ) 


OBJ= $(TEST_OBJ) $(UTIL_OBJ) $(MGMT_OBJ)  $(ASN1C_MODULES) $(E2AP_OBJ) $(E2SM_OBJ) $(XAPP_OBJ)

print-%  : ; @echo $* = $($*)

hw_unit_tests: $(OBJ)
	$(CXX) -o $@  $(OBJ) $(LIBS) $(RNIBFLAGS) $(CPPFLAGS) $(CLOGFLAGS)

install: hw_unit_tests
	install  -D hw_unit_tests  /usr/local/bin/hw_unit_tests

clean:
	-rm *.o $(E2APSRC)/*.o $(UTILSRC)/*.o $(E2SMSRC)/*.o  $(MGMTSRC)/*.o $(SRC)/*.o hw_unit_tests 
