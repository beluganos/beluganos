#! /bin/bash
# -*- coding: utf-8 -*-

# Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

CNFOPT=" --config-file /etc/beluganos/fibc.conf"
LOGOPT=" --log-config-file /etc/beluganos/fibc.log.conf"

LOGFILE=/tmp/fibc.log.trace

do_usage() {
    echo "$0 debug"
    echo "  show detail messages."
    echo ""
    echo "$0 trace"
    echo "  debug mode and output stdout to file."
    echo ""
    echo "$0 silent"
    echo "  no messages."
    echo ""
    echo "$0"
    echo "  show some messages."
    echo ""
    echo "$0 clean"
    echo "  remove log files."
}

case $1 in
    debug)
	echo "DEBUG MODE"
	ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp --verbose $CNFOPT $LOGOPT
	;;

    trace)
	echo "TRACE MODE"
	ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp --verbose $CNFOPT $LOGOPT 2>&1 | tee $LOGFILE
	;;

    silent)
	echo "SILENT MODE"
	ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp $CNFOPT > /dev/null 2>&1
	;;

    clean)
	rm -fv /tmp/fibc.log*
	;;

    help)
	do_usage
	;;

    *)
	echo "RELEASE MODE"
	ryu-manager ryu.app.ofctl_rest fabricflow.fibc.app.fibcapp $CNFOPT
	;;
esac
