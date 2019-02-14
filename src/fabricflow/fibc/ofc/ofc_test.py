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
fabricflow.fibc.ofc init
"""

import unittest

from fabricflow.fibc.ofc import ofc
from fabricflow.fibc.api import fibcapi_pb2 as pb

DP_TYPES = [
    "default",
    "generic",
    "ovs",
    "ofdpa2",
    "onsl",
]

class TestFIBCOfc(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass

    def test_get_mod(self):
        for dp_type in DP_TYPES:
            self.assertIsNotNone(ofc.get_mod(dp_type, 0, -1, None))

    def test_flow(self):
        tables = [
            -1,
            pb.FlowMod.VLAN,
            pb.FlowMod.TERM_MAC,
            pb.FlowMod.MPLS1,
            pb.FlowMod.UNICAST_ROUTING,
            pb.FlowMod.BRIDGING,
            pb.FlowMod.POLICY_ACL,
        ]

        for dp_type in DP_TYPES:
            for table in tables:
                self.assertIsNotNone(ofc.flow(dp_type, table))


    def test_group(self):
        groups = [
            -1,
            pb.GroupMod.L2_INTERFACE,
            pb.GroupMod.L3_UNICAST,
            pb.GroupMod.L3_ECMP,
            pb.GroupMod.MPLS_INTERFACE,
            pb.GroupMod.MPLS_L3_VPN,
            pb.GroupMod.MPLS_TUNNEL1,
            pb.GroupMod.MPLS_SWAP,
            pb.GroupMod.MPLS_ECMP,
            pb.GroupMod.L2_UF_INTERFACE,
        ]

        for dp_type in DP_TYPES:
            for group in groups:
                self.assertIsNotNone(ofc.group(dp_type, group))


    def test_func(self):
        func_names = [
            "pkt_out",
            "get_port_stats",
        ]

        for dp_type in DP_TYPES:
            for func_name in func_names:
                self.assertIsNotNone(getattr(ofc, func_name)(dp_type))


TESTS = [TestFIBCOfc]

if __name__ == "__main__":
    unittest.main()
