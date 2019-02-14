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
FIB Controller database
"""

class FIBCDbDpTable(object):
    """
    Table of DP ports,
    """
    def __init__(self, dpset):
        self.dpext = dict() # datapaths not in ryu-dpset.
        self.dpset = dpset
        self.dpcfg = dict()


    def add_entry(self, dpcfg):
        """
        Add Entry
        """
        self.dpcfg[dpcfg["dp_id"]] = dpcfg


    def del_entry(self, dp_id):
        """
        Delete Entry
        """
        return self.dpcfg.pop(dp_id, None)


    def show(self):
        """
        Show datas.
        """
        pass


    def dump(self, writer):
        """
        Dump datas.
        """
        import json
        json.dump(self.dpcfg, writer)


    def get_mode(self, dp_id, default=None):
        """
        Get Dp Mode
        """
        dpcfg = self.dpcfg.get(dp_id, None)
        if dpcfg is not None:
            return dpcfg["mode"]

        return default


    def add_dp(self, dpath):
        """
        Add extended datapath
        """
        dp_id = dpath.id
        if dp_id not in self.dpext:
            self.dpext[dp_id] = dpath


    def del_dp(self, dp_id):
        """
        Delete extended datapath
        """
        return self.dpext.pop(dp_id, None)


    def keys(self):
        """
        Get DpId list.
        """
        def _check_mode(dpid):
            mode = self.get_mode(dpid)
            return mode is not None and mode != "default"

        dpids = self.dpset.dps.keys()
        dpids.extend(self.dpext.keys())
        return [dpid for dpid in dpids if _check_mode(dpid)]


    def find_by_id(self, dp_id):
        """
        find (datapath, mode) by dp_id.
        """
        dp_id = int(str(dp_id))
        dpath = self.dpext.get(dp_id, None)
        if dpath is None:
            dpath = self.dpset.get(dp_id)

        if dpath is None:
            raise KeyError("{0} is not found in dpset.".format(dp_id))

        mode = self.get_mode(dp_id, "default")

        return dpath, mode


    def find_by_name(self, name):
        """
        find by name
        """
        for _, entry in self.dpcfg.items():
            if entry["name"] == name:
                return entry

        raise KeyError("{0} is not found.".format(name))

    def find_port(self, dp_id, port_id):
        """
        find port
        """
        dp_id = int(str(dp_id))
        return self.dpset.get_port(dp_id, port_id)
