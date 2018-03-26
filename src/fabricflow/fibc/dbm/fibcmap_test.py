#! /usr/bin/env python
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

import unittest

from fabricflow.fibc.dbm import fibcmap

class TestFIBMaps(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass


    def test_new_map_entry(self):
        e = fibcmap.new_map_entry(10, 1, "1.1.1.1", "eth1")
        self.assertEqual(e.port.id, 10)
        self.assertEqual(e.port.port, 1)
        self.assertEqual(e.link.re_id, "1.1.1.1")
        self.assertEqual(e.link.name, "eth1")


    def test_fibc_map(self):
        m = fibcmap.FIBMaps()

        e11 = fibcmap.new_map_entry(10, 1, "1.1.1.1", "eth1")
        e12 = fibcmap.new_map_entry(10, 2, "1.1.1.1", "eth2")
        e21 = fibcmap.new_map_entry(20, 1, "1.1.1.2", "eth1")
        e22 = fibcmap.new_map_entry(20, 2, "1.1.1.2", "eth2")

        m.insert(e11)
        m.insert(e12)
        m.insert(e21)
        m.insert(e22)

        self.assertEqual(m.find_by_port(10, 1), e11)
        self.assertEqual(m.find_by_port(10, 2), e12)
        self.assertEqual(m.find_by_port(20, 1), e21)
        self.assertEqual(m.find_by_port(20, 2), e22)

        self.assertEqual(m.find_by_link("1.1.1.1", "eth1"), e11)
        self.assertEqual(m.find_by_link("1.1.1.1", "eth2"), e12)
        self.assertEqual(m.find_by_link("1.1.1.2", "eth1"), e21)
        self.assertEqual(m.find_by_link("1.1.1.2", "eth2"), e22)


    def test_fibc_table_vm_trigger(self):
        t = fibcmap.FIBMapTable()
        t.maps.insert(fibcmap.new_map_entry(10, 1, "1.1.1.1", "eth1"))
        t.maps.insert(fibcmap.new_map_entry(10, 2, "1.1.1.1", "eth2"))
        t.maps.insert(fibcmap.new_map_entry(20, 1, "1.1.1.2", "eth1"))
        t.maps.insert(fibcmap.new_map_entry(20, 2, "1.1.1.2", "eth2"))

        t.insert_vm("1.1.1.1", "eth1")
        t.insert_vm("1.1.1.2", "eth2")

        e = t.find_by_name("1.1.1.1", "eth1")
        self.assertEqual(t.find_by_dp(10, 1), e)
        self.assertEqual(e["vm"],   fibcmap.FIBCPort(id="1.1.1.1",    port=0))
        self.assertEqual(e["dp"],   fibcmap.FIBCPort(id=10,           port=1))
        self.assertEqual(e["vs"],   fibcmap.FIBCPort(id=0 ,           port=0))
        self.assertEqual(e["name"], fibcmap.FIBCLink(re_id="1.1.1.1", name="eth1"))

        e = t.find_by_name("1.1.1.1", "eth2")
        self.assertEqual(t.find_by_dp(10, 2), None)
        self.assertEqual(e, None)

        e = t.find_by_name("1.1.1.2", "eth1")
        self.assertEqual(t.find_by_dp(20, 1), None)
        self.assertEqual(e, None)

        e = t.find_by_name("1.1.1.2", "eth2")
        self.assertEqual(t.find_by_dp(20, 2), e)
        self.assertEqual(e["vm"],   fibcmap.FIBCPort(id="1.1.1.2",    port=0))
        self.assertEqual(e["dp"],   fibcmap.FIBCPort(id=20,           port=2))
        self.assertEqual(e["vs"],   fibcmap.FIBCPort(id=0 ,           port=0))
        self.assertEqual(e["name"], fibcmap.FIBCLink(re_id="1.1.1.2", name="eth2"))


    def test_fibc_table_dp_trigger(self):
        t = fibcmap.FIBMapTable()
        t.maps.insert(fibcmap.new_map_entry(10, 1, "1.1.1.1", "eth1"))
        t.maps.insert(fibcmap.new_map_entry(10, 2, "1.1.1.1", "eth2"))
        t.maps.insert(fibcmap.new_map_entry(20, 1, "1.1.1.2", "eth1"))
        t.maps.insert(fibcmap.new_map_entry(20, 2, "1.1.1.2", "eth2"))

        t.insert_dp(10, 1)
        t.insert_dp(20, 2)

        e = t.find_by_dp(10, 1)
        self.assertEqual(t.find_by_name("1.1.1.1", "eth1"), e)
        self.assertEqual(e["vm"],   fibcmap.FIBCPort(id="1.1.1.1",    port=0))
        self.assertEqual(e["dp"],   fibcmap.FIBCPort(id=10,           port=1))
        self.assertEqual(e["vs"],   fibcmap.FIBCPort(id=0 ,           port=0))
        self.assertEqual(e["name"], fibcmap.FIBCLink(re_id="1.1.1.1", name="eth1"))

        e = t.find_by_dp(10, 2)
        self.assertEqual(t.find_by_name("1.1.1.1", "eth2"), None)
        self.assertEqual(e, None)

        e = t.find_by_dp(20, 1)
        self.assertEqual(t.find_by_name("1.1.1.2", "eth1"), None)
        self.assertEqual(e, None)

        e = t.find_by_dp(20, 2)
        self.assertEqual(t.find_by_name("1.1.1.2", "eth2"), e)
        self.assertEqual(e["vm"],   fibcmap.FIBCPort(id="1.1.1.2",    port=0))
        self.assertEqual(e["dp"],   fibcmap.FIBCPort(id=20,           port=2))
        self.assertEqual(e["vs"],   fibcmap.FIBCPort(id=0 ,           port=0))
        self.assertEqual(e["name"], fibcmap.FIBCLink(re_id="1.1.1.2", name="eth2"))


TESTS = [TestFIBMaps]

if __name__ == "__main__":
    import logging
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()
