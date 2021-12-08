ACLOCAL_AMFLAGS = -I m4

redismoduledir = @libexecdir@/$(PACKAGE)
redismodule_LTLIBRARIES = libredismodule.la
#lib_LTLIBRARIES = libredismodule.la

BASE_LDFLAGS =

libredismodule_la_SOURCES = \
	include/redismodule.h\
	src/exstrings.c

libredismodule_la_CFLAGS = \
	-std=c11 -fPIC -g -Wall -Werror -Wextra \
	-DREDISMODULE_EXPERIMENTAL_API \
	-I${includedir} -I${top_srcdir}/include


libredismodule_la_LDFLAGS = $(BASE_LDFLAGS) -avoid-version -module -shared

#pkgincludedir = ${includedir}
#pkginclude_HEADERS = include/redismodule.h

clean-local:
	rm -rf ${builddir}/libredismodule.pc

if UNIT_TEST_ENABLED
# UT
CPP_U_TEST=$(CPP_U_TEST_LATEST)
check_PROGRAMS = redismodule_ut redismodule_ut2
#TESTS = ${check_PROGRAMS}
redismodule_ut_SOURCES = \
	src/exstrings.c \
	tst/mock/include/commonStub.h \
	tst/mock/include/exstringsStub.h \
	tst/mock/include/redismodule.h  \
	tst/mock/src/commonStub.cpp \
	tst/mock/src/redismoduleStub.cpp \
	tst/src/exstrings_test.cpp \
	tst/src/main.cpp


redismodule_ut_CFLAGS = \
		-std=c11 -g -Wall \
		-fprofile-arcs -ftest-coverage \
		-D__UT__ \
		$(LIBCPPUTEST_CFLAGS) \
		-I${top_srcdir}/tst/mock/include \
		-I${includedir} \
		-I${top_srcdir}/include \
		-I${CPP_U_TEST_LATEST}/include \
		-Wall -Werror -Wextra

redismodule_ut_CXXFLAGS = \
	-std=c++11 -g -Wall \
	-fprofile-arcs -ftest-coverage \
	-D__UT__ \
	$(LIBCPPUTEST_CFLAGS) \
	-I${top_srcdir}/tst/mock/include \
	-I${includedir} \
	-I${top_srcdir}/include \
	-I${CPP_U_TEST_LATEST}/include \
	-Wall -Werror -Wextra


redismodule_ut_LDFLAGS = -Wl,-rpath=${libdir} ${UT_COVERAGE_LDFLAGS}
redismodule_ut_LDADD = -L${libdir} $(LIBCPPUTEST_LIBS) -lgcov


redismodule_ut2_SOURCES = \
	src/exstrings.c \
	tst/include/ut_helpers.hpp \
	tst/mock/include/commonStub.h \
	tst/mock/include/exstringsStub.h \
	tst/mock/include/redismodule.h  \
	tst/mock/src/commonStub.cpp \
	tst/mock/src/redismoduleNewStub.cpp \
	tst/src/exstrings_ndel_test.cpp \
	tst/src/exstrings_nget_test.cpp \
	tst/src/main.cpp \
	tst/src/ut_helpers.cpp


redismodule_ut2_CFLAGS = \
		-std=c11 -g -Wall \
		-fprofile-arcs -ftest-coverage \
		-D__UT__ \
		$(LIBCPPUTEST_CFLAGS) \
		-I${top_srcdir}/tst/mock/include \
		-I${includedir} \
		-I${top_srcdir}/include \
		-I${CPP_U_TEST_LATEST}/include \
		-Wall -Werror -Wextra

redismodule_ut2_CXXFLAGS = \
	-std=c++11 -g -Wall \
	-fprofile-arcs -ftest-coverage \
	-D__UT__ \
	$(LIBCPPUTEST_CFLAGS) \
	-I${top_srcdir}/tst/mock/include \
	-I${top_srcdir}/tst/include \
	-I${includedir} \
	-I${top_srcdir}/include \
	-I${CPP_U_TEST_LATEST}/include \
	-Wall -Werror -Wextra


redismodule_ut2_LDFLAGS = -Wl,-rpath=${libdir} ${UT_COVERAGE_LDFLAGS}
redismodule_ut2_LDADD = -L${libdir} $(LIBCPPUTEST_LIBS) -lgcov


if UNIT_TEST_MEMCHECK_ENABLED
test: ut_memcheck_test
else
test: ut_test
endif

ut_test: redismodule_ut redismodule_ut2
	./redismodule_ut
	./redismodule_ut2

ut_memcheck_test: redismodule_ut redismodule_ut2
	valgrind --error-exitcode=1 --leak-check=full ./redismodule_ut
	valgrind --error-exitcode=1 --leak-check=full ./redismodule_ut2

TESTS = run-tests.sh

else

test:
	echo 'enable ut with configure flag: --enable-unit-test and valgrind memcheck with: --enable-unit-test-memcheck'
	exit 1

endif #UNIT_TEST_ENABLED
