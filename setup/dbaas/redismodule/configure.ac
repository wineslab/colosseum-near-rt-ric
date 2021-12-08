AC_INIT([redismodule], [0.0.0], [], [], [https://gerrit.oran-osc.org/r/#/admin/projects/ric-plt/dbaas])
AC_CONFIG_AUX_DIR([build-aux])
AC_CONFIG_MACRO_DIRS([m4])
AC_CONFIG_HEADERS([config.h])


AM_INIT_AUTOMAKE([-Wall -Werror foreign subdir-objects tar-pax])

m4_ifdef([AM_PROG_AR], [AM_PROG_AR])
LT_INIT([disable-static])

# Checks for programs.
AC_PROG_CXX
AC_PROG_CC
AC_PROG_INSTALL
AC_PROG_MAKE_SET

AC_ARG_ENABLE([unit-test],
  [--disable-unit-test], [Disable unit tests]
)

AS_IF([test x"$enable_unit_test" != x"no"],
  [
    PKG_CHECK_MODULES([LIBCPPUTEST],[cpputest])
    AM_CONDITIONAL([UNIT_TEST_ENABLED], [true])
  ],
  [
    AM_CONDITIONAL([UNIT_TEST_ENABLED], [false])
  ]
)

AC_ARG_ENABLE([unit-test-memcheck],
  [AS_HELP_STRING([--disable-unit-test-memcheck], [Do not run unit tests with valgrind])]
)

AS_IF([test x"$enable_unit_test_memcheck" != x"no"],
  [
    AC_CHECK_PROG(VALGRIND_CHECK, valgrind, yes)
    AM_CONDITIONAL([UNIT_TEST_MEMCHECK_ENABLED], [true])
  ],
  [
    AM_CONDITIONAL([UNIT_TEST_MEMCHECK_ENABLED], [false])
  ]
)


# Checks for header files.

# Checks for typedefs, structures, and compiler characteristics.
AC_HEADER_STDBOOL

AC_CONFIG_FILES([Makefile])

AC_OUTPUT
