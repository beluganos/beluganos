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
FIB Controller Web API
"""

import logging
import json
from StringIO import StringIO

from ryu.base import app_manager
from ryu.app.wsgi import ControllerBase
from ryu.app.wsgi import Response
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.net import ffpacket
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.lib import fibcevt

_LOG = logging.getLogger(__name__)

class FIBCRestApp(app_manager.RyuApp):
    """
    Rest App
    """

    _EVENTS = [
        fibcevt.EventFIBCPortMap,
        fibcevt.EventFIBCPortConfig,
        fibcevt.EventFIBCDpPortConfig,
        fibcevt.EventFIBCVsPortConfig,
        fibcevt.EventFIBCFlowMod,
        fibcevt.EventFIBCGroupMod,
    ]

    def create(self, wsgi):
        """
        Setup app
        """
        wsgi.register(FIBCRestController, {"app": self})

        mapper = wsgi.mapper

        path = "/fib"

        url = path
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_all",
                       condition=dict(method=["GET"]))


        url = path + "/portmap"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_portmap",
                       condition=dict(method=["GET"]))


        url = path + "/portmap/{table}/{cmd}"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="mod_portmap",
                       condition=dict(method=["POST"]))


        url = path + "/idmap"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_idmap",
                       condition=dict(method=["GET"]))


        url = path + "/dpmap"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_dpmap",
                       condition=dict(method=["GET"]))


        url = path + "/portcfg/{name}"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="port_config",
                       condition=dict(method=["POST"]))


        url = path + "/flow/{table}"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="mod_flow",
                       condition=dict(method=["POST"]))


        url = path + "/group/{name}"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="mod_group",
                       condition=dict(method=["POST"]))


class FIBCRestController(ControllerBase):
    """
    Rest Controller
    """
    def __init__(self, req, link, data, **config):
        super(FIBCRestController, self).__init__(req, link, data, **config)
        self.app = data["app"]

    # pylint: disable=unused-argument
    @staticmethod
    def get_all(req, **kwargs):
        """
        Get All datas.
        """
        sio = StringIO()
        fibcdbm.dump(sio)
        return Response(content_type='application/json', body=sio.getvalue())


    @staticmethod
    def get_portmap(req, **kwargs):
        """
        Get All datas.
        """
        sio = StringIO()
        fibcdbm.portmap().dump(sio)
        return Response(content_type='application/json', body=sio.getvalue())


    @staticmethod
    def get_idmap(req, **kwargs):
        """
        Get All datas.
        """
        sio = StringIO()
        fibcdbm.idmap().dump(sio)
        return Response(content_type='application/json', body=sio.getvalue())


    @staticmethod
    def get_dpmap(req, **kwargs):
        """
        Get All datas.
        """
        sio = StringIO()
        fibcdbm.dps().dump(sio)
        print sio.getvalue()
        return Response(content_type='application/json', body=sio.getvalue())


    def mod_portmap(self, req, table, cmd, **kwargs):
        """
        modify portmap
        """
        dic = json.loads(req.body)
        if table == "dp":
            msg = fibcdbm.create_dp(dic)
            evt = fibcevt.EventFIBCPortMap(msg, cmd, "dp")
            self.app.send_event_to_observers(evt)

        elif table == "re":
            msg = fibcdbm.create_idmap(dic)
            evt = fibcevt.EventFIBCPortMap(msg, cmd, "idmap")
            self.app.send_event_to_observers(evt)

        elif table == "port":
            for port in fibcdbm.create_ports(dic):
                evt = fibcevt.EventFIBCPortMap(port, cmd, "port")
                self.app.send_event_to_observers(evt)

        else:
            pass


    def port_config(self, req, name, **kwargs):
        """
        simulate port config
        """
        msg = json.loads(req.body)
        if name == "vm":
            self.vm_port_config(msg)
        elif name == "vs":
            self.vs_port_config(msg)
        elif name == "dp":
            self.dp_port_config(msg)
        else:
            _LOG.error("bad name. %s", name)
            return Response(status=404)

        return Response(status=200)


    def vm_port_config(self, msg):
        """
        simulate vm port config event.
        """
        cmd = msg["cmd"]
        re_id = msg["re_id"]
        for arg in msg["args"]:
            cfg = pb.PortConfig(cmd=cmd, re_id=re_id, ifname=arg["ifname"], value=arg["port"])
            evt = fibcevt.EventFIBCPortConfig(cfg)
            self.app.send_event_to_observers(evt)


    def dp_port_config(self, msg):
        """
        simulate dp enter event
        """
        dp_id = msg["dp_id"]
        for arg in msg["args"]:
            evt = fibcevt.EventFIBCDpPortConfig(msg, dp_id, arg["port"], arg["enter"])
            self.app.send_event_to_observers(evt)


    def vs_port_config(self, msg):
        """
        simulate ffpacket received.
        """
        print msg
        re_id = msg["re_id"]
        vs_id = msg["vs_id"]
        for arg in msg["args"]:
            pkt = ffpacket.FFPacket(re_id, arg["ifname"])
            evt = fibcevt.EventFIBCVsPortConfig(pkt, vs_id, arg["port"])
            self.app.send_event_to_observers(evt)


    def mod_flow(self, req, table, **kwargs):
        """
        simulate flow mod event
        """
        msg = json.loads(req.body)
        if table == "vlan":
            self._mod_flow_vlan(msg)

        elif table == "termmac":
            self._mod_flow_termmac(msg)

        elif table == "mpls1":
            self._mod_flow_mpls1(msg)

        elif table == "unicast":
            self._mod_flow_unicast(msg)

        else:
            _LOG.error("bad flow mod %s %s", table, msg)
            return Response(status=404)

        return Response(status=200)


    def _mod_flow_vlan(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            match = pb.VLANFlow.Match(in_port=arg["port"],
                                      vid=arg.get("vid", 0),
                                      vid_mask=arg.get("mask", 0))

            actions = []
            if arg["vrf"] != 0:
                actions.append(pb.VLANFlow.Action(name="SET_VRF", value=arg["vrf"]))
            if "new_vid" in arg:
                actions.append(pb.VLANFlow.Action(name="PUSH_VLAN", value=arg["new_vid"]))

            flow = pb.VLANFlow(match=match, actions=actions, goto_table=20)
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="VLAN", vlan=flow)
            evt = fibcevt.EventFIBCFlowMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_flow_termmac(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            match = pb.TerminationMacFlow.Match(eth_type=arg["eth_type"], eth_dst=arg["eth_dst"])
            flow = pb.TerminationMacFlow(match=match, actions=[], goto_table=arg["goto"])
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="TERM_MAC", term_mac=flow)
            evt = fibcevt.EventFIBCFlowMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_flow_mpls1(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            match = pb.MPLSFlow.Match(bos=arg["bos"], label=arg["label"])
            actions = []
            for name in ("SET_VRF", "POP_LABEL"):
                if name in arg:
                    actions.append(pb.MPLSFlow.Action(name=name, value=arg[name]))
            flow = pb.MPLSFlow(match=match,
                               actions=actions,
                               g_type=arg.get("g_type", "UNSPEC"),
                               g_id=arg.get("g_id", 0),
                               goto_table=arg["goto"])
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="MPLS1", mpls1=flow)
            evt = fibcevt.EventFIBCFlowMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_flow_unicast(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            match = pb.UnicastRoutingFlow.Match(ip_dst=arg["dst"], vrf=arg["vrf"])
            flow = pb.UnicastRoutingFlow(match=match,
                                         action=None,
                                         g_type=arg.get("g_type", "UNSPEC"),
                                         g_id=arg.get("g_id", 0))
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="UNICAST_ROUTING", unicast=flow)
            evt = fibcevt.EventFIBCFlowMod(mod)

            self.app.send_event_to_observers(evt)


    def mod_group(self, req, name, **kwargs):
        """
        simulare group mod event
        """
        msg = json.loads(req.body)
        if name == "l2_interface":
            self._mod_group_l2_interface(msg)

        elif name == "l3_unicast":
            self._mod_group_l3_unicast(msg)

        elif name == "mpls_interface":
            self._mod_group_mpls_interface(msg)

        elif name == "mpls_l2_vpn":
            pass

        elif name == "mpls_l3_vpn":
            self._mod_group_l3_vpn(msg)

        elif name == "mpls_tun1":
            self._mod_group_mpls_tun1(msg)

        elif name == "mpls_tun2":
            pass

        elif name == "mpls_swap":
            self._mod_group_mpls_swap(msg)

        else:
            _LOG.error("bad group mod %s %s", name, msg)
            return Response(status=404)

        return Response(status=200)


    def _mod_group_l2_interface(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.L2InterfaceGroup(port_id=arg["port"], vlan_vid=arg["vlan"])
            mod = pb.GroupMod(cmd=cmd, g_type="L2_INTERFACE", re_id=re_id, l2_iface=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_l3_unicast(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.L3UnicastGroup(ne_id=arg["ne_id"],
                                    port_id=arg["port"],
                                    vlan_vid=arg["vlan"],
                                    eth_dst=arg["eth_dst"],
                                    eth_src=arg["eth_src"])
            mod = pb.GroupMod(cmd=cmd, g_type="L3_UNICAST", re_id=re_id, l3_unicast=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_mpls_interface(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSInterfaceGroup(ne_id=arg["ne_id"],
                                        port_id=arg["port"],
                                        vlan_vid=arg["vlan"],
                                        eth_dst=arg["eth_dst"],
                                        eth_src=arg["eth_src"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_INTERFACE", re_id=re_id, mpls_iface=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_l3_vpn(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSLabelGroup(label=arg["label"],
                                    out_label=arg["out_label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_L3_VPN", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_mpls_tun1(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSLabelGroup(label=arg["label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_TUNNEL1", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_mpls_swap(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSLabelGroup(label=arg["label"],
                                    out_label=arg["out_label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_SWAP", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)
