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
FIBC mod factory Tests
"""

import unittest
from fabricflow.fibc.ofc import generic

class TestGeneric(unittest.TestCase):
    """
    TestCase for generic module.
    """

    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_fix_port_stats_names(self):
        """
        _fix_port_stats_names
        """
        data_1234 = {'rx_packets': 10,
                     'tx_packets': 11,
                     'rx_bytes': 12,
                     'tx_bytes': 13,
                     'rx_dropped': 14,
                     'tx_dropped': 15,
                     'rx_errors': 16,
                     'tx_errors': 17,
                     'rx_frame_err': 18,
                     'rx_over_err': 19,
                     'rx_crc_err': 20,
                     'collisions': 21,
                     'duration_sec': 22,
                     'duration_nsec': 23,
                     "port_no": 24}
        data = {"1234": [data_1234]}


        res = generic._fix_port_stats_names(data) # pylint: disable=protected-access
        res_1234 = res["1234"][0]
        self.assertEqual(res_1234.get("ifInUcastPkts", -1), data_1234["rx_packets"])
        self.assertEqual(res_1234.get("ifOutUcastPkts", -1), data_1234["tx_packets"])
        self.assertEqual(res_1234.get("ifInNUcastPkts", -1), 0)
        self.assertEqual(res_1234.get("ifOutNUcastPkts", -1), 0)
        self.assertEqual(res_1234.get("ifInOctets", -1), data_1234["rx_bytes"])
        self.assertEqual(res_1234.get("ifOutOctets", -1), data_1234["tx_bytes"])
        self.assertEqual(res_1234.get("ifInDiscards", -1), data_1234["rx_dropped"])
        self.assertEqual(res_1234.get("ifOutDiscards", -1), data_1234["tx_dropped"])
        self.assertEqual(res_1234.get("ifInErrors", -1), data_1234["rx_errors"])
        self.assertEqual(res_1234.get("ifOutErrors", -1), data_1234["tx_errors"])
        self.assertEqual(res_1234.get("etherStatsJabbers", -1), data_1234["rx_frame_err"])
        self.assertEqual(res_1234.get("etherStatsOversizePkts", -1), data_1234["rx_over_err"])
        self.assertEqual(res_1234.get("etherStatsCRCAlignErrors", -1), data_1234["rx_crc_err"])
        self.assertEqual(res_1234.get("etherStatsCollisions", -1), data_1234["collisions"])


TESTS = [TestGeneric]

if __name__ == "__main__":
    unittest.main()
