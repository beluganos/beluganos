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

import grpc
from nlaapi import nlaapi_pb2 as api

def _main():
    channel = grpc.insecure_channel('127.0.0.1:50052')
    stub = api.NLAApiStub(channel)
    for vpn in stub.GetVpns(api.GetVpnsRequest()):
        print vpn

if __name__ == "__main__":
    _main()
