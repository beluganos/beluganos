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
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.api import fibcapi as api


class TestVLANFlow(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        m = pb.VLANFlow.Match(in_port = 1, vid = 10)
        a = [
            pb.VLANFlow.Action(name="SET_VLAN_VID", value=10),
            pb.VLANFlow.Action(name="SET_VRF", value=1)
        ]
        flow = pb.VLANFlow(match=m, actions=a, goto_table=20)
        mod = pb.FlowMod(cmd="ADD", table="VLAN", vlan=flow)

        b = mod.SerializeToString()
        m = api.parse_flow_mod(b)

        self.assertEqual(str(mod), str(m))


class TestTermMACFlow(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        match = pb.TerminationMacFlow.Match(eth_type=0x0800, eth_dst="11:22:33:44:55:66")
        actions = []
        flow = pb.TerminationMacFlow(match=match, actions=actions, goto_table=30)
        mod = pb.FlowMod(cmd="ADD", table="TERM_MAC", term_mac=flow)

        b = mod.SerializeToString()
        m = api.parse_flow_mod(b)

        self.assertEqual(str(mod), str(m))


class TestMPLS1Flow(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        match = pb.MPLSFlow.Match(bos=1, label=10017)
        actions = [
            pb.MPLSFlow.Action(name="POP_LABEL", value=0)
        ]
        flow = pb.MPLSFlow(match=match, actions=actions, goto_table=60, g_type=pb.GroupMod.MPLS_INTERFACE, g_id=1)
        mod = pb.FlowMod(cmd="ADD", table="MPLS1", mpls1=flow)
        # print mod


class TestL2InterfaceGroup(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        group = pb.L2InterfaceGroup(port_id=1, vlan_vid=10)
        mod = pb.GroupMod(cmd="ADD", g_type="L2_INTERFACE", re_id="1.1.1.1", l2_iface=group)
        # print mod


class TestMplsLabelGroup(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        group = pb.MPLSLabelGroup(dst_id=16, new_label=10016, ne_id=1, g_type="MPLS_INTERFACE")
        mod = pb.GroupMod(cmd="ADD", g_type="MPLS_SWAP", re_id="1.1.1.1", mpls_label=group)
        # print mod


class TestPortConfig(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new(self):
        msg = pb.PortConfig(cmd="ADD", re_id="1.1.1.1", ifname="ethX", port_id=10)
        # print msg


class TestMultipart(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new_request_port(self):
        dp_id = 0x123456789a
        port_no = 1
        m = api.new_ff_multipart_request_port(dp_id, port_no, "")
        self.assertEqual(m.mp_type, pb.FFMultipart.PORT)

        b = m.SerializeToString()
        m = api.parse_ff_multipart_request(b)
        self.assertEqual(m.port, pb.FFMultipart.PortRequest(port_no=1))


class TestFFPortMod(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_new_ff_port_mod(self):
        pm = api.new_ff_port_mod(0x1234, 10, pb.PortStatus.UP)
        self.assertEqual(pm.dp_id, 0x1234)
        self.assertEqual(pm.port_no, 10)
        self.assertEqual(pm.status, pb.PortStatus.UP)

        pm = api.new_ff_port_mod(0x1234, 10, "UP")
        self.assertEqual(pm.dp_id, 0x1234)
        self.assertEqual(pm.port_no, 10)
        self.assertEqual(pm.status, pb.PortStatus.UP)


TESTS = [
    TestPortConfig,
    TestVLANFlow,
    TestTermMACFlow,
    TestMPLS1Flow,
    TestL2InterfaceGroup,
    TestMplsLabelGroup,
    TestMultipart,
    TestFFPortMod,
]

if __name__ == "__main__":
    unittest.main()
                                            
