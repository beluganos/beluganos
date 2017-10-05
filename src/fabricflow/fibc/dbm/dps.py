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
        self.dpset = dpset
        self.dpmap = dict()


    def _create(self, dpmap, dpset):
        """
        Create table.
        """
        self.dpset = dpset
        self.dpmap = dpmap


    def add(self, dpath):
        """
        Add Entry
        """
        self.dpmap[dpath["dp_id"]] = dpath


    def delete_by_dp_id(self, dp_id):
        """
        Delete Entry
        """
        return self.dpmap.pop(dp_id, None)


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
        json.dump(self.dpmap, writer)


    def find_port(self, dp_id, port_id):
        """
        Find Port by dpid/port
        """
        try:
            return self.dpset.get_port(dp_id, port_id)

        except Exception as expt:
            raise KeyError(str(expt))


    def get_mode(self, dp_id, default=None):
        """
        Get Dp Mode
        """
        dpcfg = self.dpmap.get(dp_id, None)
        if dpcfg is not None:
            return dpcfg["mode"]

        return default


    def find_by_id(self, dp_id):
        """
        find (datapath, mode) by dp_id.
        """
        dpath = self.dpset.get(dp_id)
        if dpath is None:
            raise KeyError("{0} is not found in dpset.".format(dp_id))
        return dpath


    def find_by_name(self, name):
        """
        find by name
        """
        for _, entry in self.dpmap.items():
            if entry["name"] == name:
                return entry

        raise KeyError("{0} is not found.".format(name))
