# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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
FIB Controller database (Portmap)
"""

import logging
from collections import namedtuple
from fabricflow.fibc.dbm.portmap import FIBCPort
from fabricflow.fibc.dbm.portmap import FIBCLink
from fabricflow.fibc.dbm.portmap import FIBCPortEntry

_LOG = logging.getLogger(__name__)

FIBMap = namedtuple('FIBMap', ("port", "link")) # port:FIBCPort, link:FIBCLink


def new_map_entry(dp_id, dp_port, re_id, name):
    """
    returns instance of FIBMap
    """
    port = FIBCPort(id=dp_id, port=dp_port)
    link = FIBCLink(re_id=re_id, name=name)
    return FIBMap(port=port, link=link)


class FIBMaps(object):
    """
    port-ifname mapping.
    """

    def __init__(self):
        self.ports = dict()
        self.links = dict()


    def insert(self, entry):
        """
        Insert entry(FIBMap)
        """
        if entry.port in self.ports or entry.link in self.links:
            raise KeyError

        self.ports[entry.port] = entry
        self.links[entry.link] = entry


    def find_by_port(self, dp_id, dp_port):
        """
        Select entry by port.
        """
        port = FIBCPort(id=dp_id, port=dp_port)
        return self.ports[port]


    def find_by_link(self, re_id, name):
        """
        Select entry by name.
        """
        link = FIBCLink(re_id=re_id, name=name)
        return self.links[link]


class FIBMapTable(object):
    """
    Table of FIBCPortEntry instances.
    """

    def __init__(self):
        self.entries = dict()
        self.maps = FIBMaps()

    def add(self, entry):
        """
        Insert entry(FIBCPortEntry)
        """
        key = entry["name"]
        self.entries[key] = entry
        return entry


    def insert_vm(self, re_id, name):
        """
        Insert entry(vm trigger)
        """
        mentry = self.maps.find_by_link(re_id, name)
        entry = FIBCPortEntry.new(
            dp_id=mentry.port.id,
            port=mentry.port.port,
            re_id=re_id,
            name=name,
        )

        return self.add(entry)


    def insert_dp(self, dp_id, dp_port):
        """
        Insert entry(dp trigger)
        """
        mentry = self.maps.find_by_port(dp_id, dp_port)
        entry = FIBCPortEntry.new(
            dp_id=dp_id,
            port=dp_port,
            re_id=mentry.link.re_id,
            name=mentry.link.name,
        )

        return self.add(entry)


    def delete_by_name(self, re_id, name):
        """
        Delete entry(vm trigger)
        """
        key = FIBCLink(re_id=re_id, name=name)
        if key in self.entries:
            return self.entries.pop(key, None)

        return None


    def find_by_key(self, key):
        """
        find entry by key(FIBCLink)
        """
        return self.entries.get(key, None)


    def find_by_name(self, re_id, name):
        """
        find entry by vm(re_id, ifname)
        """
        key = FIBCLink(re_id=re_id, name=name)
        return self.find_by_key(key)


    def find_by_dp(self, dp_id, dp_port):
        """
        find entry by dp(dp_id, dp_port)
        """
        key = FIBCPort(id=dp_id, port=dp_port)
        for entry in self.entries.values():
            if entry["dp"] == key:
                return entry

        return None


    def find_by_vs(self, vs_id, vs_port):
        """
        find entry by vs(vs_id, vs_port)
        """
        key = FIBCPort(id=vs_id, port=vs_port)
        for entry in self.entries.values():
            if entry["vs"] == key:
                return entry

        return None


    def list_by_link_ref(self, link):
        """
        Find port by link
        """
        entries = []
        for entry in self.entries.values():
            ref = entry.get("link", None)
            if ref is not None and ref == link:
                entries.append(entry)

        return sorted(entries)


    def slaves(self, entry):
        """
        list slaves(FIBCLink instances)
        """
        slaves = entry.get("slaves", None)
        if slaves is None:
            return []

        return [self.find_by_key(slave) for slave in slaves]


    def master(self, slave):
        """
        find maser(FINCLink instance)
        """
        for entry in self.entries.values():
            slaves = entry.get("slaves", None)
            if slaves and slave["name"] in slaves:
                return entry

        return None


    def subs(self, entry, exclude=True):
        """
        list sub interface entries(FIBCLink instances)
        """
        subs = self.list_by_link_ref(entry["name"])
        if not subs and not exclude:
            return [entry]

        entries = []
        for sub in subs:
            entries.extend(self.subs(sub), False)

        return entries


    def parent(self, entry):
        """
        find top of parent(FIBCLinl instance)
        """
        link = entry.get("link", None)
        if link is None:
            return entry

        return self.parent(self.find_by_key(link))
