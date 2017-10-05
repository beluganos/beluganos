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

import logging

_LOG = logging.getLogger(__name__)

class FIBCDbIdMapEntry(dict):
    """
    Entry for idmap table.
    """
    def __init__(self, re_id, dp_id, vm_status=False, dp_status=False):
        super(FIBCDbIdMapEntry, self).__init__(self)
        self["re_id"] = re_id
        self["dp_id"] = dp_id
        self["vm_status"] = vm_status
        self["dp_status"] = dp_status


    def associated(self):
        """
        check vm/dp associated.
        """
        return self["vm_status"] and self["dp_status"]


    def update_vm_status(self, status):
        """
        update vm status
        """
        old_status = self.associated()
        self["vm_status"] = status
        return old_status != self.associated()


    def update_dp_status(self, status):
        """
        update dp status.
        """
        old_status = self.associated()
        self["dp_status"] = status
        return old_status != self.associated()


class FIBCDbIdMapTable(object):
    """
    Table for dp_id/re_id conversion.
    """
    def __init__(self):
        self.entries = dict()
        self.dps = dict()


    def add(self, re_id, dp_id):
        """
        Add entry
        """
        if (re_id not in self.entries) and (dp_id not in self.dps):
            entry = FIBCDbIdMapEntry(re_id, dp_id)
            self.entries[re_id] = entry
            self.dps[dp_id] = entry


    def delete_by_re_id(self, re_id):
        """
        Delete entry
        """
        entry = self.entries.pop(re_id, None)
        if entry:
            self.dps.pop(entry["dp_id"])

        return entry


    def show(self):
        """
        Show datas.
        """
        for entry in self.entries.values():
            _LOG.debug("%s", entry)


    def dump(self, writer):
        """
        Dump datas.
        """
        import json
        json.dump(self.entries.values(), writer)


    def find_by_re_id(self, re_id):
        """
        Find entry by re_id.
        """
        return self.entries[re_id]


    def find_by_dp_id(self, dp_id):
        """
        Find entry by dp_id
        """
        return self.dps[dp_id]
