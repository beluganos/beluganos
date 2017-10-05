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

"""
FIBC Logging module
"""

# pylint: disable=too-few-public-methods
class FIBCLog(object):
    """
    Logger settings
    """
    TRACE_LOG = False
    DUMP_PKT = False
    DUMP_MSG = False


def set_trace(enable):
    """
    enable/disable trace log
    """
    FIBCLog.TRACE_LOG = enable


def trace():
    """
    trace log status.
    """
    return FIBCLog.TRACE_LOG


def set_dump_pkt(enable):
    """
    enable/disable dump log
    """
    FIBCLog.DUMP_PKT = enable


def dump_pkt():
    """
    dump log status.
    """
    return FIBCLog.DUMP_PKT


def set_dump_msg(enable):
    """
    enable/disable message log
    """
    FIBCLog.DUMP_MSG = enable


def dump_msg():
    """
    message log status.
    """
    return FIBCLog.DUMP_MSG
