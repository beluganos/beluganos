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

AC_ARG_WITH([opennsl],
            [AS_HELP_STRING([--with-opennsl], [enable OpenNSL. @<:@default=no@:>@])],
            [],
            [with_opennsl=no])

AM_CONDITIONAL([ENABLE_OPENNSL], [test "$with_opennsl" = "yes"])

#
# Makefiles
#
AC_CONFIG_FILES([Makefile
		 etc/Makefile
		 etc/installer/Makefile
		 src/Makefile
		 src/goryu/Makefile
		 src/goryu/ryulib/Makefile
		 src/goryu/ofproto/Makefile
		 src/goryu/encoding/Makefile
		 src/gonsl/Makefile
		 src/gonsl/lib/Makefile
		 src/gonsl/api/Makefile
		 src/fabricflow/Makefile
		 src/fabricflow/util/Makefile
		 src/fabricflow/util/net/Makefile
		 src/fabricflow/util/container/Makefile
		 src/fabricflow/util/container/mqueue/Makefile
		 src/fabricflow/fibc/Makefile
		 src/fabricflow/fibc/net/Makefile
		 src/fabricflow/fibc/api/Makefile
		 src/fabricflow/fibc/app/Makefile
		 src/fabricflow/fibc/dbm/Makefile
		 src/fabricflow/fibc/lib/Makefile
		 src/fabricflow/fibc/ofc/Makefile
		 src/fabricflow/fibs/Makefile
		 src/fabricflow/ribc/Makefile
		 src/fabricflow/ribc/ribctl/Makefile
		 src/fabricflow/ribp/Makefile
		 src/fabricflow/ribp/api/Makefile
		 src/fabricflow/ribp/ribpkt/Makefile
		 src/fabricflow/ribs/Makefile
		 src/fabricflow/ribs/api/Makefile
		 src/fabricflow/ribs/api/ribsapi/Makefile
		 src/fabricflow/ribs/pkgs/Makefile
		 src/fabricflow/ribs/pkgs/ribscfg/Makefile
		 src/fabricflow/ribs/pkgs/ribsdbm/Makefile
		 src/fabricflow/ribs/pkgs/ribssrv/Makefile
		 src/fabricflow/ribt/Makefile
		 src/fabricflow/ribt/api/Makefile
		 src/fabricflow/ribt/api/ribtapi/Makefile
		 src/fabricflow/ribn/Makefile
		 src/fabricflow/ffctl/Makefile
		 src/gonla/Makefile
		 src/gonla/nlactl/Makefile
		 src/gonla/nlamsg/Makefile
		 src/gonla/nlamsg/nlalink/Makefile
		 src/gonla/nlalib/Makefile
		 src/gonla/nlaapi/Makefile
		 src/gonla/nladbm/Makefile
		 src/gonla/nlasvc/Makefile
		 src/govsw/Makefile
		 src/govsw/api/Makefile
		 src/govsw/api/vswapi/Makefile
		 src/govsw/pkgs/Makefile
		 src/govsw/pkgs/govsw/Makefile
	])
AC_OUTPUT
