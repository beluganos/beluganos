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
OFDPA Extended Actions
"""

from ryu.ofproto import ofproto_v1_3_parser as ofparser

OFDPA_EXPERIMENTER_ID = 0x00001018

class OfdpaAction(ofparser.OFPActionExperimenter):
    """
    OFDPA Extended Action
    """
    _experimenter = OFDPA_EXPERIMENTER_ID
    type = 10
    len = 8
    fmt = "!HHI"

    def __init__(self, port):
        super(OfdpaAction, self).__init__(experimenter=self._experimenter)
        self.port = port


    @classmethod
    def parse(cls, buf, offset):
        """
        Parse binary message
        """
        pass


    def serialize(self, buf, offset):
        """
        Serialize to binary message.
        """
        from ryu.lib.pack_utils import msg_pack_into
        msg_pack_into(self.fmt, buf,
                      offset, self.type, self.len, self.port)


def _generate(ofpp_name):
    import sys
    ofpp = sys.modules[ofpp_name]

    def _add_attr(key, val):
        val.__module__ = ofpp.__name__  # Necessary for stringify stuff
        setattr(ofpp, key, val)

    _add_attr("OfdpaAction", OfdpaAction)


def init():
    """
    Activate OFDPA Extention
    """
    import logging
    logger = logging.getLogger(__name__)
    _generate("ryu.ofproto.ofproto_v1_3_parser")
    logger.info("%s loaded.", __name__)
