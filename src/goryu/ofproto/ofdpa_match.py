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
OFDPA Extended Matches
"""

from ryu.lib import type_desc
from ryu.ofproto import oxm_fields

OFDPA_OXM_VRF = 1
OFDPA_OXM_TRAFFIC_CLASS = 2
OFDPA_OXM_COLOR = 3
OFDPA_OXM_DEI = 4
OFDPA_OXM_QOS_INDEX = 5
OFDPA_OXM_LMEP_ID = 6
OFDPA_OXM_MPLS_TTL = 7
OFDPA_OXM_MPLS_L2_PORT = 8
OFDPA_OXM_L3_IN_PORT = 9
OFDPA_OXM_OVID = 10
OFDPA_OXM_MPLS_DATA_FIRST_NIBBLE = 11
OFDPA_OXM_MPLS_ACH_CHANNEL = 12
OFDPA_OXM_MPLS_NEXT_LABEL_IS_GAL = 13
OFDPA_OXM_OAM_Y1731_MDL = 14
OFDPA_OXM_OAM_Y1731_OPCODE = 15
OFDPA_OXM_COLOR_ACTIONS_INDEX = 16
OFDPA_OXM_PROTECTION_INDEX = 21
OFDPA_OXM_ETH_SUB_TYPE = 22
OFDPA_OXM_MPLS_TYPE = 23
OFDPA_OXM_ALLOW_VLAN_TRANSLATION = 24

OFDPA_EXPERIMENTER_ID = 0x00001018


# pylint: disable=protected-access
# pylint: disable=too-few-public-methods
class OfdpaExperimenter(oxm_fields._Experimenter):
    """
    OFDPA Experimenter
    """
    experimenter_id = OFDPA_EXPERIMENTER_ID

# pylint: disable=bad-whitespace
_OXM_TYPES = [
    OfdpaExperimenter("vrf",                     OFDPA_OXM_VRF,                    type_desc.Int2),
    OfdpaExperimenter("traffic_class",           OFDPA_OXM_TRAFFIC_CLASS,          type_desc.Int1),
    OfdpaExperimenter("color",                   OFDPA_OXM_COLOR,                  type_desc.Int1),
    OfdpaExperimenter("dei",                     OFDPA_OXM_DEI,                    type_desc.Int1),
    OfdpaExperimenter("qos_index",               OFDPA_OXM_QOS_INDEX,              type_desc.Int1),
    OfdpaExperimenter("lemp_id",                 OFDPA_OXM_LMEP_ID,                type_desc.Int4),
    OfdpaExperimenter("mpls_ttl",                OFDPA_OXM_MPLS_TTL,               type_desc.Int1),
    OfdpaExperimenter("mpls_l2_port",            OFDPA_OXM_MPLS_L2_PORT,           type_desc.Int4),
    OfdpaExperimenter("l3_in_port",              OFDPA_OXM_L3_IN_PORT,             type_desc.Int4),
    OfdpaExperimenter("ovid",                    OFDPA_OXM_OVID,                   type_desc.Int2),
    OfdpaExperimenter("mpls_data_first_nibble",  OFDPA_OXM_MPLS_DATA_FIRST_NIBBLE, type_desc.Int1),
    OfdpaExperimenter("mpls_ach_channel",        OFDPA_OXM_MPLS_ACH_CHANNEL,       type_desc.Int2),
    OfdpaExperimenter("mpls_next_label_is_gal",  OFDPA_OXM_MPLS_NEXT_LABEL_IS_GAL, type_desc.Int1),
    OfdpaExperimenter("y1731_mdl",               OFDPA_OXM_OAM_Y1731_MDL,          type_desc.Int1),
    OfdpaExperimenter("y1731_opcode",            OFDPA_OXM_OAM_Y1731_OPCODE,       type_desc.Int1),
    OfdpaExperimenter("color_action_index",      OFDPA_OXM_COLOR_ACTIONS_INDEX,    type_desc.Int4),
    OfdpaExperimenter("protection_index",        OFDPA_OXM_PROTECTION_INDEX,       type_desc.Int1),
    OfdpaExperimenter("eth_sub_type",            OFDPA_OXM_ETH_SUB_TYPE,           type_desc.Int1),
    OfdpaExperimenter("mpls_type",               OFDPA_OXM_MPLS_TYPE,              type_desc.Int2),
    OfdpaExperimenter("allow_vlan_translation",  OFDPA_OXM_ALLOW_VLAN_TRANSLATION, type_desc.Int1)
]

def init():
    """
    Activate OFDPA Extention
    """
    import logging
    from ryu.ofproto import ofproto_v1_3 as ofproto
    logger = logging.getLogger(__name__)
    ofproto.oxm_types += _OXM_TYPES
    oxm_fields.generate("ryu.ofproto.ofproto_v1_3")
    logger.info("%s loaded.", __name__)
