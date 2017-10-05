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
Group Mod functions
"""

def is_bucket_needed(dpath, cmd):
    """
    judge buckets are needed.
    """
    ofp = dpath.ofproto
    return cmd != ofp.OFPGC_DELETE


def group_mod(group_id, gtype, buckets):
    """
    group mod
    """
    return dict(
        type=gtype,
        group_id=group_id,
        buckets=buckets() if callable(buckets) else buckets,
    )


def clear_all(dpath):
    """
    Clear all groups.
    """
    ofp = dpath.ofproto

    def group_clear(gtype):
        """
        Clear groups by type.
        """
        # return {"group_id": dpath.ofproto.OFPG_ALL, 'type': gtype}
        return dpath.ofproto_parser.OFPGroupMod(
            dpath, ofp.OFPGC_DELETE, gtype, ofp.OFPG_ALL, [])

    gtypes = [ofp.OFPGT_ALL, ofp.OFPGT_SELECT, ofp.OFPGT_INDIRECT, ofp.OFPGT_FF]
    return [group_clear(gtype) for gtype in gtypes]
