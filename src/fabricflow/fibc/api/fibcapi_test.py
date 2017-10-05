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
from goryu.ofproto import offlow
from fabricflow.fibc.api import fibcapi


class TestVLANFlow(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_1(self):
        pass


class TestApi(unittest.TestCase):
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_priority(self):
        datas = [
            ("1.0.0.0/24", 16400+offlow.PRIORITY_BAND*24, 16400),
            ("1.0.0.0/24", 16405+offlow.PRIORITY_BAND*24, 16405),
            ("1.0.0.0/16", 16405+offlow.PRIORITY_BAND*16, 16405),
        ]
        for ip, pri, base in datas:
            self.assertEqual(offlow.priority_for_ipaddr(ip, base), pri)


TESTS = [TestApi]

if __name__ == "__main__":
    unittest.main()
