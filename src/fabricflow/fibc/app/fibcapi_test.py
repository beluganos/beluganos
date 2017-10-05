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


import mock
import unittest
from fabricflow.fibc.app import fibcapi


class TestFIBCApiClients(unittest.TestCase):
    def setUp(self):
        pass


    def tearDown(self):
        pass


    def test_register_unregister(self):
        sock = mock.Mock()
        clients = fibcapi.FIBCApiClients()

        # empty
        self.assertEqual(len(clients.clients), 0)

        # 1 client
        clients.register("1.1.1.1", sock, "192.168.1.1")
        self.assertEqual(len(clients.clients), 1)
        clients.unregister("1.1.1.1")
        self.assertEqual(len(clients.clients), 0)
        sock.close.assert_called_once_with()


    def test_get_sock(self):
        sock = mock.Mock()
        clients = fibcapi.FIBCApiClients()

        # not found
        with self.assertRaises(KeyError):
            clients.get_sock("1.1.1.1")

        # 1 client
        clients.register("1.1.1.1", sock, "192.168.1.1")
        ret = clients.get_sock("1.1.1.1")
        self.assertEqual(ret, sock)

        with self.assertRaises(KeyError):
            clients.get_sock("1.1.1.2")

        # remoce client
        clients.unregister("1.1.1.1")
        with self.assertRaises(KeyError):
            clients.get_sock("1.1.1.1")


TESTS = [TestFIBCApiClients]

if __name__ == "__main__":
    unittest_man()
