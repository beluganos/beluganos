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
FIB Controller database (Portmap)
"""

import logging
from collections import namedtuple

_LOG = logging.getLogger(__name__)

FIBCPort = namedtuple('FIBCPort', ("id", "port"))
FIBCLink = namedtuple('FIBCLink', ("re_id", "name"))

class FIBCPortEntry(dict):
    """
    Entry of Portmap table.
    """

    def __init__(self, name, vm, dp, vs, **kwargs):
        super(FIBCPortEntry, self).__init__()
        self["name"] = name
        self["vm"] = vm
        self["vs"] = vs
        self["dp"] = dp
        self["link"] = kwargs.get("link", None)
        self["slaves"] = kwargs.get("slaves", None)
        self["dpenter"] = kwargs.get("dpenter", None)


    @classmethod
    def new(cls, **kwargs):
        """
        return new instance.
        """
        dp_id = kwargs["dp_id"]
        dp_port = kwargs["port"]
        vm_id = kwargs["re_id"]
        vm_port = kwargs.get("vm_port", 0)
        vs_id = kwargs.get("vs_id", 0)
        vs_port = kwargs.get("vs_port", 0)
        link = kwargs.get("link", None)
        if link:
            link = FIBCLink(vm_id, link)
        else:
            link = None
        slaves = kwargs.get("slaves", None)
        if slaves is not None:
            slaves = [FIBCLink(vm_id, slave) for slave in slaves]
        return cls(
            dp=FIBCPort(dp_id, dp_port),
            vm=FIBCPort(vm_id, vm_port),
            vs=FIBCPort(vs_id, vs_port),
            name=FIBCLink(vm_id, kwargs["name"]),
            link=link,
            slaves=slaves,
            dpenter=kwargs.get("dpenter", False),
        )


    def update_vs(self, vs_id, port_id):
        """
        Replace vs port.
        """
        if self["vs"].id == vs_id and self["vs"].port == port_id:
            return False

        self["vs"] = FIBCPort(vs_id, port_id)
        return True


    def update_vm(self, port_id):
        """
        Replace vm port_id
        """
        self["vm"] = FIBCPort(self["vm"].id, port_id)
        return self


    def is_datapath_ready(self):
        """
        Check if vs,dp port associated.
        """
        return self["dpenter"] and self["vs"].port != 0


    def is_vm_ready(self):
        """
        Check if VM port associated
        """
        return self["vm"].port != 0


class FIBCDbPortMapTable(object):
    """
    Table for port conversion.
    """

    def __init__(self):
        self.ports = dict()


    def add(self, port):
        """
        Add Port
        """
        key = port["name"]
        self.ports[key] = port


    def delete_by_name_key(self, key):
        """
        Delete port by key(FIBCLink)
        """
        if key in self.ports:
            return self.ports.pop(key, None)


    def delete_by_name(self, re_id, name):
        """
        Delete port by re_id/ifname
        """
        key = FIBCLink(re_id=re_id, name=name)
        return self.delete_by_name_key(key)


    def show(self):
        """
        Show datas.
        """
        for port in self.ports.values():
            _LOG.debug("%s", port)


    def dump(self, writer):
        """
        Dump Ports
        """
        import json
        json.dump(self.ports.values(), writer)


    def find_by_name_key(self, key):
        """
        Find port by name(FIBCLink)
        """
        return self.ports[key]


    def find_by_name(self, re_id, name):
        """
        Find port by re_id and ifname.
        """
        key = FIBCLink(re_id=re_id, name=name)
        return self.find_by_name_key(key)


    def find_by_key(self, kind, pid, port_id):
        """
        find port by key
        """
        key = FIBCPort(id=pid, port=port_id)
        for port in self.ports.values():
            if key == port[kind]:
                return port

        raise KeyError("{0} not found.".format(key))


    def list_by_link_ref(self, key):
        """
        Find port by link.
        """
        ports = []
        for port in self.ports.values():
            link = port.get("link", None)
            if link is not None and key == link:
                ports.append(port)

        return sorted(ports)


    def find_by_vm(self, re_id, port_id):
        """
        Find port (key=FIBCPort())
        """
        return self.find_by_key("vm", re_id, port_id)


    def find_by_vs(self, vs_id, port_id):
        """
        Find port by vs_id and vs port_id
        """
        return self.find_by_key("vs", vs_id, port_id)


    def find_by_dp(self, dp_id, port_id):
        """
        Find port by dp_id and dp port_id
        """
        return self.find_by_key("dp", dp_id, port_id)


    def slave_ports(self, port):
        """
        Get slave port list.
        """
        slaves = port["slaves"]
        if slaves is None:
            return []

        return [self.find_by_name_key(slave) for slave in slaves]


    def master_port(self, port):
        """
        Get master port
        """
        for master in self.ports.values():
            slaves = master.get("slaves", None)
            if slaves is None:
                continue

            if port["name"] in slaves:
                return master

        raise KeyError("Not slave device. {0}".format(port))


    def upper_ports(self, port, exclude=True):
        """
        Get Upper port
        """
        upper_ports = self.list_by_link_ref(port["name"])
        if not upper_ports and not exclude:
            return [port]

        ports = []
        for upper_port in upper_ports:
            ports.extend(self.upper_ports(upper_port, False))

        return ports


    def lower_port(self, port):
        """
        Get lower port.
        """
        link = port.get("link", None)
        if link is None:
            return port

        port = self.find_by_name_key(link)
        return self.lower_port(port)
