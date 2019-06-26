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
from mock import Mock
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibccnv

_RE_ID = "1.1.1.1"
_DP_ID = 100
_IFNAME = "ethX"
_VID = 10
_VID_MASK = 0x1fff
_ETH_SRC_MAC = "77:88:99:aa:bb:cc"
_ETH_DST_MAC = "11:22:33:44:55:66"
_ETH_DST_MASK = "01:02:03:04:05:06"
_ETH_DST = _ETH_DST_MAC+_ETH_DST_MASK

class TestConvFlow(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_vlan(self):
        flow = pb.VLANFlow(
            match=pb.VLANFlow.Match(
                in_port=1,
                vid=_VID,
                vid_mask=_VID_MASK,
            ),
            actions=[],
            goto_table=pb.FlowMod.TERM_MAC
        )
        mod = pb.FlowMod(
            cmd="ADD",
            table=pb.FlowMod.VLAN,
            re_id=_RE_ID,
            vlan=flow
        )

        p = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p
        portmap.lower_port.return_value = p

        #exec
        fibccnv.conv_flow(mod, portmap)

        #check
        self.assertEqual(mod.cmd, pb.FlowMod.ADD)
        self.assertEqual(mod.table, pb.FlowMod.VLAN)
        self.assertEqual(mod.re_id, _RE_ID)
        self.assertEqual(mod.vlan.match.in_port, 2)
        self.assertEqual(mod.vlan.match.vid, _VID)
        self.assertEqual(mod.vlan.match.vid_mask, _VID_MASK)
        self.assertEqual(len(mod.vlan.actions), 0)
        self.assertEqual(mod.vlan.goto_table, pb.FlowMod.TERM_MAC)


    def test_term_mac(self):
        flow = pb.TerminationMacFlow(
            match=pb.TerminationMacFlow.Match(
                in_port=1,
                eth_type=0x0800,
                eth_dst=_ETH_DST,
                vlan_vid=_VID
            ),
            actions=[],
            goto_table=pb.FlowMod.UNICAST_ROUTING
        )
        mod = pb.FlowMod(
            cmd="ADD",
            table=pb.FlowMod.TERM_MAC,
            re_id=_RE_ID,
            term_mac=flow
        )

        p = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p
        portmap.lower_port.return_value = p

        #exec
        fibccnv.conv_flow(mod, portmap)

        #check
        self.assertEqual(mod.cmd, pb.FlowMod.ADD)
        self.assertEqual(mod.table, pb.FlowMod.TERM_MAC)
        self.assertEqual(mod.re_id, _RE_ID)
        self.assertEqual(mod.term_mac.match.in_port, 2)
        self.assertEqual(mod.term_mac.match.eth_type, 0x0800)
        self.assertEqual(mod.term_mac.match.eth_dst, _ETH_DST)
        self.assertEqual(mod.term_mac.match.vlan_vid, _VID)
        self.assertEqual(len(mod.term_mac.actions), 0)
        self.assertEqual(mod.term_mac.goto_table, pb.FlowMod.UNICAST_ROUTING)


    def test_bridging(self):
        flow = pb.BridgingFlow(
            action=pb.BridgingFlow.Action(
                name=pb.PolicyACLFlow.Action.OUTPUT,
                value=1,
            )
        )
        mod = pb.FlowMod(
            cmd="ADD",
            table=pb.FlowMod.BRIDGING,
            re_id=_RE_ID,
            bridging=flow,
        )
        p = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p
        portmap.lower_port.return_value = p

        # exec
        fibccnv.conv_flow(mod, portmap)
        # print mod

        # check
        self.assertEqual(mod.cmd, pb.FlowMod.ADD)
        self.assertEqual(mod.table, pb.FlowMod.BRIDGING)
        self.assertEqual(mod.re_id, _RE_ID)
        self.assertEqual(mod.bridging.action.name, pb.PolicyACLFlow.Action.OUTPUT)
        self.assertEqual(mod.bridging.action.value, 2)


class TestConvGroup(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_conv_l2_interface(self):
        group = pb.L2InterfaceGroup(
            port_id = 1,
            vlan_vid = _VID,
            vlan_translation = True
        )
        mod = pb.GroupMod(re_id=_RE_ID, g_type="L2_INTERFACE", l2_iface=group)

        p = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p
        portmap.lower_port.return_value = p

        # exec
        fibccnv.conv_group(mod, portmap)
        # check
        self.assertEqual(mod.l2_iface.port_id, 2)
        portmap.find_by_vm.assert_called_once_with(re_id=_RE_ID, port_id=1)
        portmap.lower_port.assert_called_once_with(p)


    def test_l3_unicast(self):
        group = pb.L3UnicastGroup(
            ne_id=1010,
            port_id=1,
            vlan_vid=_VID,
            eth_dst=_ETH_DST_MAC,
            eth_src=_ETH_SRC_MAC,
        )
        mod = pb.GroupMod(re_id=_RE_ID, g_type="L3_UNICAST", l3_unicast=group)

        p = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p
        portmap.lower_port.return_value = p

        # exec
        fibccnv.conv_group(mod, portmap)
        # check
        self.assertEqual(mod.l3_unicast.port_id, 2)
        portmap.find_by_vm.assert_called_once_with(re_id=_RE_ID, port_id=1)
        portmap.lower_port.assert_called_once_with(p)


    def test_mpls_iface(self):
        group = pb.MPLSInterfaceGroup(
            ne_id=1010,
            port_id=1,
            vlan_vid=_VID,
            eth_dst=_ETH_DST_MAC,
            eth_src=_ETH_SRC_MAC,
        )
        mod = pb.GroupMod(re_id=_RE_ID, g_type="MPLS_INTERFACE", mpls_iface=group)

        p2 = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=2, dp_id=_DP_ID, re_id=_RE_ID)
        p3 = fibcdbm.FIBCPortEntry.new(name=_IFNAME, port=3, dp_id=_DP_ID, re_id=_RE_ID)
        portmap = Mock(spec=fibcdbm.FIBCDbPortMapTable)
        portmap.find_by_vm.return_value = p2
        portmap.lower_port.return_value = p3

        # exec
        fibccnv.conv_group(mod, portmap)
        # check
        self.assertEqual(mod.mpls_iface.port_id, 3)
        portmap.find_by_vm.assert_called_once_with(re_id=_RE_ID, port_id=1)
        portmap.lower_port.assert_called_once_with(p2)


TESTS = [TestConvFlow, TestConvGroup]

if __name__ == "__main__":
    unittest.main()
