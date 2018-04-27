# -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.69])
AC_INIT([GoFlow], [0.1], [admin@example.net])
AM_INIT_AUTOMAKE([foreign])

# Checks for programs.
AC_PROG_INSTALL

# Checks for libraries.

# Checks for header files.

# Checks for typedefs, structures, and compiler characteristics.

# Checks for library functions.

FF_HOME=`pwd`
AC_SUBST([FF_HOME])

AC_ARG_WITH([go-prefix],
            [AS_HELP_STRING([--with-go-prefix],
	    [GOPATH @<:@default=~/go@:>@])],
	    [],
	    [with_go_prefix=~/go])

GOPREFIX="$with_go_prefix"
AC_SUBST([GOPREFIX])

#
# Makefiles
#
AC_CONFIG_FILES([Makefile
                 src/Makefile
		 src/goryu/Makefile
		 src/goryu/ryulib/Makefile
		 src/goryu/ofproto/Makefile
		 src/goryu/encoding/Makefile
		 src/fabricflow/Makefile
		 src/fabricflow/util/Makefile
		 src/fabricflow/util/net/Makefile
		 src/fabricflow/fibc/Makefile
		 src/fabricflow/fibc/net/Makefile
		 src/fabricflow/fibc/api/Makefile
		 src/fabricflow/fibc/app/Makefile
		 src/fabricflow/fibc/dbm/Makefile
		 src/fabricflow/fibc/lib/Makefile
		 src/fabricflow/fibc/ofc/Makefile
		 src/fabricflow/ribc/Makefile
		 src/fabricflow/ribc/ribctl/Makefile
		 src/fabricflow/ribp/Makefile
		 src/fabricflow/ribp/api/Makefile
		 src/fabricflow/ribp/ribpkt/Makefile
		 src/fabricflow/ribs/Makefile
		 src/fabricflow/ribs/ribsapi/Makefile
		 src/fabricflow/ribs/ribsmsg/Makefile
		 src/fabricflow/ribs/ribsyn/Makefile
		 src/gonla/Makefile
		 src/gonla/nlactl/Makefile
		 src/gonla/nlamsg/Makefile
		 src/gonla/nlamsg/nlalink/Makefile
		 src/gonla/nlalib/Makefile
		 src/gonla/nlaapi/Makefile
		 src/gonla/nladbm/Makefile
		 src/gonla/nlasvc/Makefile
	])
AC_OUTPUT
