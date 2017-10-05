#! /usr/bin/env python
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
beluganos/ control command.
"""

import logging
import subprocess

_LOG = logging.getLogger()
_LOG_DEBUG_FORMAT = "%(asctime)s %(levelname)s %(message)s [%(filename)s:%(lineno)d]"
_LOG_DEFAULT_FORMAT = "%(levelname)s %(message)s"

_OFC = "tcp:127.0.0.1:6633"

def init_ovs(bridge):
    """
    Initializxe OVS.
    """
    cmds = [
        "sudo ovs-vsctl add-br {0}".format(bridge),
        "sudo ovs-vsctl set bridge {0} protocols=OpenFlow13".format(bridge),
        "sudo ovs-vsctl set-controller {0} {1}".format(bridge, _OFC)
    ]
    for cmd in cmds:
        ret = subprocess.call(cmd.split(" "))
        _LOG.debug("[%s/init] %d %s", bridge, ret, cmd)


def clear_ovs(bridge):
    """
    Clear OVS bridge.
    """
    cmd = "sudo ovs-vsctl del-br {0}".format(bridge)
    ret = subprocess.call(cmd.split(" "))
    _LOG.debug("[%s/clear] %d %s", bridge, ret, cmd)


def show_status():
    """
    Show all status.
    """
    cmds = [
        "systemctl status --no-pager fibcd",
        "lxc list",
        "sudo ovs-vsctl show",
    ]
    for cmd in cmds:
        subprocess.call(cmd.split(" "))


def show_container(container):
    """
    Show container status.
    """
    cmd = "lxc info {0}".format(container)
    subprocess.call(cmd.split(" "))


def start_container(name):
    """
    Start container
    """
    cmd = "lxc start {0}".format(name)
    subprocess.check_call(cmd.split(" "))
    _LOG.info("[%s] stated.", name)


def stop_container(name):
    """
    Stop container
    """
    cmd = "lxc stop {0}".format(name)
    subprocess.call(cmd.split(" "))
    _LOG.info("[%s] stopped.", name)


def get_device_list(profile):
    """
    Get device name list from lxc profile.
    """
    cmd = "lxc profile device list {0}".format(profile)
    ret = subprocess.check_output(cmd.split(" "))
    lines = ret.split("\n")
    def _line_to_name(line):
        return line.split(":", 1)[0].strip()

    return [_line_to_name(line) for line in lines if len(line) > 0]


def get_device_property(profile, device, name):
    """
    Get device property from profile.
    """
    cmd = "lxc profile device get {0} {1} {2}".format(profile, device, name)
    ret = subprocess.check_output(cmd.split(" "))
    return ret.strip()


def add_ovs_port(bridge, port):
    """
    Add port to OVS bridge.
    """
    cmd = "sudo ovs-vsctl add-port {0} {1}".format(bridge, port)
    return subprocess.call(cmd.split(" "))


def del_ovs_port(bridge, port):
    """
    Remove port from OVS bridge.
    """
    cmd = "sudo ovs-vsctl del-port {0} {1}".format(bridge, port)
    return subprocess.call(cmd.split(" "))


def add_device_to_ovs(profile, device, bridge):
    """
    Add device in profile to OVS bridge.
    """
    nictype = get_device_property(profile, device, "nictype")
    if nictype == "p2p":
        host_name = get_device_property(profile, device, "host_name")
        add_ovs_port(bridge, host_name)
        _LOG.info("%s.%s added to %s", profile, device, bridge)


def del_device_from_ovs(profile, device, bridge):
    """
    Remove device in profile from OVS bridge.
    """
    nictype = get_device_property(profile, device, "nictype")
    if nictype == "p2p":
        host_name = get_device_property(profile, device, "host_name")
        del_ovs_port(bridge, host_name)
        _LOG.info("[%s] %s removed from to %s", profile, host_name, bridge)


def exec_lxc_cmd(name, cmd):
    """
    execute command on container.
    """
    cmdline = "lxc exec {0} {1}".format(name, cmd)
    subprocess.call(cmdline.split(" "))


def run_fibcd():
    fibc_pkg = "fabricflow.fibc.app.fibcapp"
    confpath = "/etc/beluganos/fibc.conf"
    logpath = "/etc/beluganos/fibc.log.conf"

    opts = "--config-file {0} --log-config-file {1}".format(confpath, logpath)
    opts += " --verbose"

    cmd = "ryu-manager ryu.app.ofctl_rest {0} {1}".format(fibc_pkg, opts)

    subprocess.call(cmd.split(" "))


def _getopts():
    cmds = ["run", "start", "stop", "add", "del", "status", "con", "init", "clear"]
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("cmd", nargs=None, choices=cmds, type=str)
    parser.add_argument("container", nargs="?", type=str)
    parser.add_argument("-p", "--profile", nargs="?", type=str, default="")
    parser.add_argument("-b", "--bridge", nargs="?", type=str, default="dp0")
    parser.add_argument("-e", "--exclude", nargs="+", type=str, default=["eth0"])
    parser.add_argument("-v", "--verbose", action="store_true", default=False)

    args = parser.parse_args()
    if not args.profile:
        args.profile = args.container

    return args


def _init_log(verbose):
    if verbose:
        logging.basicConfig(level=logging.DEBUG, format=_LOG_DEBUG_FORMAT)
    else:
        logging.basicConfig(level=logging.INFO, format=_LOG_DEFAULT_FORMAT)


def _main():
    opts = _getopts()

    _init_log(opts.verbose)

    if opts.cmd == "add":
        start_container(opts.container)
        for device in get_device_list(opts.profile):
            if device not in opts.exclude:
                add_device_to_ovs(opts.profile, device, opts.bridge)


    elif opts.cmd == "del":
        for device in get_device_list(opts.profile):
            if device not in opts.exclude:
                del_device_from_ovs(opts.profile, device, opts.bridge)
        stop_container(opts.container)

    elif opts.cmd == "run":
        run_fibcd()

    elif opts.cmd == "start":
        cmd = "sudo systemctl start fibcd"
        subprocess.call(cmd.split(" "))

    elif opts.cmd == "stop":
        cmd = "sudo systemctl stop fibcd"
        subprocess.call(cmd.split(" "))

    elif opts.cmd == "status":
        if opts.container:
            show_container(opts.container)
        else:
            show_status()

    elif opts.cmd == "con":
        exec_lxc_cmd(opts.container, "bash")

    elif opts.cmd == "init":
        init_ovs(opts.bridge)


    elif opts.cmd == "clear":
        clear_ovs(opts.bridge)


if __name__ == "__main__":
    _main()
