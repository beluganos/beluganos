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
from fabricflow.fibc.app import fibcptm

_DPID = 10
_RTID = "1.1.1.1"

class TestFIBCPtm(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_get_ready_ports(self):
        _ports = [
            # not ready
            dict(name="10/eth0", port=0, vs_port=0, vm_port=0),
            dict(name="10/eth1", port=1, vs_port=0, vm_port=0),
            dict(name="10/eth2", port=0, vs_port=1, vm_port=0),
            dict(name="10/eth3", port=0, vs_port=0, vm_port=1),
            dict(name="10/eth4", port=0, vs_port=1, vm_port=1),
            dict(name="10/eth5", port=1, vs_port=0, vm_port=1),
            dict(name="10/eth6", port=1, vs_port=1, vm_port=0),
            dict(name="10/eth6", port=1, vs_port=1, vm_port=1),
            dict(name="10/eth6", port=1, vs_port=0, vm_port=0, dpenter=True),
            dict(name="10/eth6", port=1, vs_port=1, vm_port=0, dpenter=True),
            dict(name="10/eth6", port=1, vs_port=0, vm_port=1, dpenter=True),
            # ready
            dict(name="10/eth8", port=1, vs_port=1, vm_port=1, dpenter=True),
            dict(name="10/eth7", port=0, vs_port=1, vm_port=1, dpenter=True),
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=_DPID, re_id=_RTID, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        for port in ports[:-2]:
            ready_ports = fibcptm.get_ready_ports(table, port)
            self.assertEqual(ready_ports, [])

        for port in ports[-2:]:
            ready_ports = fibcptm.get_ready_ports(table, port)
            self.assertEqual(len(ready_ports), 1)
            self.assertEqual(ready_ports[0], port)


    def test_get_ready_ports_vlan(self):
        _ports = [
            # not ready
            dict(name="10/eth1",    port=1, vs_port=1, dpenter=False),  # 0.DOWN
            dict(name="10/eth1.10", port=1, vm_port=0, link="10/eth1"), # 1.DOWN

            dict(name="10/eth2",    port=2, vs_port=0, dpenter=True),   # 2.DOWN
            dict(name="10/eth2.10", port=2, vm_port=2, link="10/eth2"), # 3.UP

            dict(name="10/eth3",    port=3, vs_port=3, dpenter=True),   # 4.UP
            dict(name="10/eth3.10", port=3, vm_port=0, link="10/eth3"), # 5.DOWN

            # ready
            dict(name="10/eth4",    port=3, vs_port=3, dpenter=True),   # 6.UP
            dict(name="10/eth4.10", port=3, vm_port=3, link="10/eth4"), # 7,UP

            dict(name="10/eth5",    port=4, vs_port=4, dpenter=True),   # 8.UP
            dict(name="10/eth5.10", port=4, vm_port=0, link="10/eth5"), # 9.DOWN
            dict(name="10/eth5.20", port=4, vm_port=4, link="10/eth5"), #10.UP

            # ready and not ready
            dict(name="10/eth6",    port=5, vs_port=5, dpenter=True),   #11.UP
            dict(name="10/eth7.10", port=5, vm_port=5, link="10/eth6"), #12.UP
            dict(name="10/eth7.20", port=5, vm_port=6, link="10/eth6"), #13.UP
        ]
        ports = [fibcdbm.FIBCPortEntry.new(dp_id=_DPID, re_id=_RTID, **p) for p in _ports]
        table = fibcdbm.FIBCDbPortMapTable()
        for port in ports:
            table.add(port)

        # ports[0..5] -> []
        for port in ports[:6]:
            ready_ports = fibcptm.get_ready_ports(table, port)
            self.assertEqual(ready_ports, [])

        # ports[6,7] -> ports[7]
        ready_ports = fibcptm.get_ready_ports(table, ports[6])
        self.assertEqual(len(ready_ports), 1)
        self.assertEqual(ready_ports[0], ports[7])

        ready_ports = fibcptm.get_ready_ports(table, ports[7])
        self.assertEqual(len(ready_ports), 1)
        self.assertEqual(ready_ports[0], ports[7])

        # ports[8] -> ports[10]
        ready_ports = fibcptm.get_ready_ports(table, ports[8])
        self.assertEqual(len(ready_ports), 1)
        self.assertEqual(ready_ports[0], ports[10])

        # ports[9] -> []
        ready_ports = fibcptm.get_ready_ports(table, ports[9])
        self.assertEqual(len(ready_ports), 0)

        # ports[10] -> ports[18]
        ready_ports = fibcptm.get_ready_ports(table, ports[10])
        self.assertEqual(len(ready_ports), 1)
        self.assertEqual(ready_ports[0], ports[10])

        # ports[11] -> ports[12,13]
        ready_ports = fibcptm.get_ready_ports(table, ports[11])
        self.assertEqual(len(ready_ports), 2)
        self.assertEqual(ready_ports[0], ports[12])
        self.assertEqual(ready_ports[1], ports[13])


TESTS = [TestFIBCPtm]

if __name__ == "__main__":
    unittest.main()
