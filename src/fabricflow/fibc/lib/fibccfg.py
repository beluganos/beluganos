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
FIB Controller configuration
"""

_KEY_ROUTERS = "routers"
_KEY_DPATHS = "datapaths"
_CFG_PATTERN = "fibc*.yml"

class PortConfig(dict):
    """
    Port config
    """

    def __init__(self, *args, **kwargs):
        super(PortConfig, self).__init__(*args, **kwargs)
        self.routers = list()
        self.dpaths = list()


    def load(self, stream):
        """
        Load config and append to fields.
        """
        import yaml
        datas = yaml.load(stream)
        self.routers.extend(datas.get(_KEY_ROUTERS, []))
        self.dpaths.extend(datas.get(_KEY_DPATHS, []))
        return self


def load_dir(dirpath):
    """
    Load config files in directory.
    """
    import glob
    import os.path
    cfg = PortConfig()
    for path in glob.glob(os.path.join(dirpath, _CFG_PATTERN)):
        with open(path, "r") as stream:
            cfg.load(stream)
    return cfg


def print_router(router):
    """
    Print router
    """
    print "Desc: '{desc}', REID: '{re_id}', DP: '{datapath}'".format(**router)
    for port in router["ports"]:
        print "  Port: {0}".format(port)


def print_datapath(dpath):
    """
    Print datapath.
    """
    print "Name: '{name}', DpId: {dp_id}/0x{dp_id:x}, Mode: '{mode}'".format(**dpath)


def _main():
    cfg = load_dir(".")

    for router in cfg.routers:
        print_router(router)

    for dpath in cfg.dpaths:
        print_datapath(dpath)


if __name__ == "__main__":
    _main()
