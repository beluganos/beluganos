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

import unittest

from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.dbm import portmap

class TestFIBCDbPortMapTable(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_create_port(self):
        dpid = 100
        reid = "1.1.1.1"
        name = "ethX"

        # exec
        ret = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name=name, port=1)

        # check
        expect = dict(
            dp=portmap.FIBCPort(dpid, 1),
            vm=portmap.FIBCPort(reid, 0),
            vs=portmap.FIBCPort(0, 0),
            name=portmap.FIBCLink(reid, name),
            link=None,
            slaves=None,
            dpenter=False,
        )
        self.assertEqual(ret, expect)

    def test_create_port_vlan(self):
        dpid = 100
        reid = "1.1.1.1"
        name = "ethX"
        link = "ethY"

        # exec
        ret = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name=name, port=1, link=link)

        # check
        expect = dict(
            dp=portmap.FIBCPort(dpid, 1),
            vm=portmap.FIBCPort(reid, 0),
            vs=portmap.FIBCPort(0, 0),
            name=portmap.FIBCLink(reid, name),
            link=portmap.FIBCLink(reid, link),
            slaves=None,
            dpenter=False,
        )
        self.assertEqual(ret, expect)

    def test_create_port_bond(self):
        dpid = 100
        reid = "1.1.1.1"
        name = "ethX"
        slaves = ["ethA", "ethB"]
 
        # exec
        ret = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name=name, port=1, slaves=slaves)

        # check
        expect = dict(
            dp=portmap.FIBCPort(dpid, 1),
            vm=portmap.FIBCPort(reid, 0),
            vs=portmap.FIBCPort(0, 0),
            name=portmap.FIBCLink(reid, name),
            link=None,
            dpenter=False,
            slaves=[
                portmap.FIBCLink(reid, slaves[0]),
                portmap.FIBCLink(reid, slaves[1]),
            ],
        )
        self.assertEqual(ret, expect)


    def test_find_by_name(self):
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1", port=1),
            dict(name="20/eth1", port=2),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        ret = table.find_by_name(reid, "10/eth1")
        self.assertEqual(ret, ports[0])

        ret = table.find_by_name(reid, "20/eth1")
        self.assertEqual(ret, ports[1])

        with self.assertRaises(KeyError):
            table.find_by_name(reid, "20/eth0")


    def test_slave_ports(self):
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1",  port=1),
            dict(name="10/eth2",  port=2),
            dict(name="10/eth3",  port=3),
            dict(name="10/eth3.10",port=0, link="10/eth3"),
            dict(name="20/bond1", port=100, slaves=["10/eth1", "10/eth2"]),
            dict(name="30/vlan1", port=0, link="20/bond1"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # not bond devce.
        port = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, name="10/eth3", port=3)
        ret = table.slave_ports(port)
        self.assertEqual(ret, [])

        # vlan tagged, not bond devce.
        port = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, name="10/eth3.10", port=0)
        ret = table.slave_ports(port)
        self.assertEqual(ret, [])

        # bond device
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="20/bond1", port=100, slaves=["10/eth1", "10/eth2"])
        ret = table.slave_ports(port)
        self.assertEqual(ret, [ports[0], ports[1]])

        # vlan tagged bond device.
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="30/vlan1", port=0, link="20/bond1")
        ret = table.slave_ports(port)
        self.assertEqual(ret, [])


    def test_master_port(self):
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1",  port=1),
            dict(name="10/eth2",  port=2),
            dict(name="10/eth3",  port=3),
            dict(name="20/bond1", port=100, slaves=["10/eth1", "10/eth2"]),
            dict(name="30/vlan1", port=0, link="20/bond1"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # 10/eth1 -> 20/bond1
        port = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[0])
        ret = table.master_port(port)
        self.assertEqual(ret, ports[3])

        # 10/eth2 -> 20/bond1
        port = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[1])
        ret = table.master_port(port)
        self.assertEqual(ret, ports[3])

        # 10/eth3 -> KeyError
        port = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[2])
        with self.assertRaises(KeyError):
            table.master_port(port)


    def test_lower_port(self):
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1", port=1),
            dict(name="10/eth2", port=2),
            dict(name="20/eth1", port=3, link="10/eth1"),
            dict(name="20/eth2", port=4, link="10/eth2"),
            dict(name="30/eth1", port=5, link="20/eth1"),
            dict(name="30/eth2", port=6, link="20/eth2"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # linked to '10/eth1'
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="30/eth1", port=6, link="20/eth1")
        ret = table.lower_port(port)
        self.assertEqual(ret, ports[0])

        # linked to '10/eth2'
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="30/eth2", port=6, link="20/eth2")
        ret = table.lower_port(port)
        self.assertEqual(ret, ports[1])

        # 2-linked device
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="30/eth0", port=0, link="20/eth0")
        with self.assertRaises(KeyError):
            table.lower_port(port)

        # not linked device.
        port = fibcdbm.FIBCPortEntry.new(
            dp_id=dpid, re_id=reid, name="10/eth1", port=6)
        ret = table.lower_port(port)
        self.assertEqual(ret, port)


    def test_list_by_link_ref(self):
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1",    port=1),
            dict(name="20/eth1",    port=2),
            dict(name="20/eth2",    port=3),
            dict(name="10/eth1.10", port=0, link="10/eth1"),
            dict(name="20/eth2.10", port=0, link="20/eth2"),
            dict(name="20/eth2.20", port=0, link="20/eth2"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # 10/eth1 -> 10/eth1.10
        link = portmap.FIBCLink(re_id=reid, name="10/eth1")
        ret = table.list_by_link_ref(link)
        self.assertEqual(ret, ports[3:4])

        # 20/eth2 -> 20/eth2.10, 20/eth2.20
        link = portmap.FIBCLink(re_id=reid, name="20/eth2")
        ret = table.list_by_link_ref(link)
        self.assertEqual(ret, ports[4:6])

        # 20/eth1 -> []
        link = portmap.FIBCLink(re_id=reid, name="20/eth1")
        ret = table.list_by_link_ref(link)
        self.assertEqual(ret, [])


    def test_upper_ports_single_level(self):
        """
        - 10/eth1 *
        """
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1",     port=1),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # 10/eth1 -> [10/eth1]
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[0])
        ret = table.upper_ports(link)
        self.assertEqual(ret, [])

        # 10/eth2 -> [10/eth2]
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, name="10/eth2", port=1)
        ret = table.upper_ports(link)
        self.assertEqual(ret, [])


    def test_upper_ports_multi_level(self):
        """
        - 10/eth2 - 10/eth2.10 *
        - 10/eth3 - 10/eth3.10 - 10/eth3.110 *
        """
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth2",     port=2),
            dict(name="10/eth2.10",  port=0, link="10/eth2"),

            dict(name="10/eth3",     port=2),
            dict(name="10/eth3.10",  port=0, link="10/eth3"),
            dict(name="10/eth3.110", port=0, link="10/eth3.10"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # 10/eth2 -> [10/eth2.10]
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[0])
        ret = table.upper_ports(link)
        self.assertEqual(ret, ports[1:2])


        # 10/eth -> [10/eth3.110]
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[2])
        ret = table.upper_ports(link)
        self.assertEqual(ret, ports[4:5])


    def test_upper_ports_multi_upper(self):
        """
        - 10/eth1
          +- 10/eth1.10 *
          +- 10/eth1.20 *
        - 10/eth2
          +- 10/eth2.10
             +- 10/eth2.110 *
             +- 10/eth2.210 *
          +- 10/eth2.20 *
        """
        dpid = 100
        reid = "1.1.1.1"
        _ports = [
            dict(name="10/eth1",     port=1),
            dict(name="10/eth1.10",  port=0, link="10/eth1"),
            dict(name="10/eth1.20",  port=0, link="10/eth1"),

            dict(name="10/eth2",     port=1),
            dict(name="10/eth2.10",  port=0, link="10/eth2"),
            dict(name="10/eth2.110", port=0, link="10/eth2.10"),
            dict(name="10/eth2.210", port=0, link="10/eth2.10"),
            dict(name="10/eth2.20",  port=0, link="10/eth2"),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # 10/eth1 -> 10/eth1.10, eth1.20
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[0])
        ret = table.upper_ports(link)
        self.assertEqual(ret, ports[1:3])

        # 10/eth2 -> eth2.110, eth2.210, 10/eth2.20
        link = fibcdbm.FIBCPortEntry.new(dp_id=dpid, re_id=reid, **_ports[3])
        ret = table.upper_ports(link)
        self.assertEqual(ret, ports[5:8])


TESTS = [TestFIBCDbPortMapTable]

if __name__ == "__main__":
    import logging
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()
