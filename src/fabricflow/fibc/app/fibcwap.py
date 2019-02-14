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
from ryu.lib import ofctl_v1_3 as ofctl
from ryu.controller import ofp_event
from ryu.controller import handler
from fabricflow.fibc.api import fibcapi_pb2 as pb
from fabricflow.fibc.net import ffpacket
from fabricflow.fibc.dbm import fibcdbm
from fabricflow.fibc.dbm.fibcdbm import FIBCPortEntry
from fabricflow.fibc.lib import fibcevt
from fabricflow.fibc.ofc import ofc

_LOG = logging.getLogger(__name__)

class FIBCRestApp(app_manager.RyuApp):
    """
    Rest App
    """

    _EVENTS = [
        fibcevt.EventFIBCPortMap,
        fibcevt.EventFIBCVmConfig,
        fibcevt.EventFIBCPortConfig,
        fibcevt.EventFIBCDpPortConfig,
        fibcevt.EventFIBCVsPortConfig,
        fibcevt.EventFIBCFlowMod,
        fibcevt.EventFIBCGroupMod,
    ]

    def __init__(self, *args, **kwargs):
        super(FIBCRestApp, self).__init__(*args, **kwargs)
        self.waiters = {}

    def create(self, wsgi):
        """
        Setup app
        """
        wsgi.register(FIBCRestController, {"app": self, "waiters": self.waiters})

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


        url = path + "/vmcfg"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="vm_config",
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


        url = path + "/dps"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_dplist",
                       condition=dict(method=["GET"]))


        url = path + "/stats/port/{dp_id}"
        mapper.connect("fib", url,
                       controller=FIBCRestController,
                       action="get_port_stats",
                       condition=dict(method=["GET"]))


    @handler.set_ev_cls(ofp_event.EventOFPPortStatsReply, handler.MAIN_DISPATCHER) # pylint: disable=no-member
    def stats_reply_ofp_handler(self, evt):
        """
        Unlock Waiter's Event.
        """
        msg = evt.msg
        dpath = msg.datapath

        if dpath.id not in self.waiters:
            return
        if msg.xid not in self.waiters[dpath.id]:
            return
        lock, msgs = self.waiters[dpath.id][msg.xid]
        msgs.append(msg)

        if msg.flags & dpath.ofproto.OFPMPF_REPLY_MORE:
            return

        del self.waiters[dpath.id][msg.xid]
        if not self.waiters[dpath.id]:
            del self.waiters[dpath.id]

        lock.set()


    @handler.set_ev_cls(fibcevt.EventFIBCMultipartReply, handler.MAIN_DISPATCHER)
    def stats_reply_fibc_handler(self, evt):
        """
        Unlock Waiter's Event.
        """
        dpath = evt.dp

        if dpath.id not in self.waiters:
            return
        if evt.xid not in self.waiters[dpath.id]:
            return

        lock, msgs = self.waiters[dpath.id][evt.xid]
        msgs.append(evt.msg)

        del self.waiters[dpath.id][evt.xid]
        if not self.waiters[dpath.id]:
            del self.waiters[dpath.id]

        lock.set()


class FIBCRestController(ControllerBase):
    """
    Rest Controller
    """
    def __init__(self, req, link, data, **config):
        super(FIBCRestController, self).__init__(req, link, data, **config)
        self.app = data["app"]
        self.waiters = data["waiters"]


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
            evt = fibcevt.EventFIBCPortMap(dic, cmd, "dp")
            self.app.send_event_to_observers(evt)

        elif table == "re":
            msg = dict(re_id=dic["re_id"], dp_id=dic["dp_id"])
            evt = fibcevt.EventFIBCPortMap(msg, cmd, "idmap")
            self.app.send_event_to_observers(evt)

        elif table == "port":
            for port in dic["ports"]:
                entry = FIBCPortEntry.new(re_id=dic["re_id"], dp_id=dic["dp_id"], **port)
                evt = fibcevt.EventFIBCPortMap(entry, cmd, "port")
                self.app.send_event_to_observers(evt)

        else:
            pass


    def vm_config(self, req, **kwargs):
        """
        simulate vm config
        """
        msg = json.loads(req.body)
        print msg
        re_id = msg["re_id"]
        enter = msg.get("cmd", "ENTER") == "ENTER"
        evt = fibcevt.EventFIBCVmConfig(pb.Hello(re_id=re_id), enter)
        self.app.send_event_to_observers(evt)

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
            cfg = pb.PortConfig(
                cmd=cmd,
                re_id=re_id,
                ifname=arg["ifname"],
                port_id=arg["port"],
                link=arg.get("link", ""),
                status=arg.get("status", "NOP"),
                dp_port=arg.get("dp_port", 0))
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

        elif table == "acl":
            self._mod_flow_acl(msg)

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
            match = pb.UnicastRoutingFlow.Match(
                ip_dst=arg["dst"],
                vrf=arg["vrf"],
                origin=arg.get("origin", "ROUTE"))
            flow = pb.UnicastRoutingFlow(match=match,
                                         action=None,
                                         g_type=arg.get("g_type", "UNSPEC"),
                                         g_id=arg.get("g_id", 0))
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="UNICAST_ROUTING", unicast=flow)
            evt = fibcevt.EventFIBCFlowMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_flow_acl(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            match = pb.PolicyACLFlow.Match(
                ip_dst=arg["dst"],
                vrf=arg.get("vrf", 0))
            action = pb.PolicyACLFlow.Action(
                name="OUTPUT",
                value=0)
            flow = pb.PolicyACLFlow(match=match, action=action)
            mod = pb.FlowMod(cmd=cmd, re_id=re_id, table="POLICY_ACL", acl=flow)
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
            grp = pb.L2InterfaceGroup(
                port_id=arg["port"],
                vlan_vid=arg["vlan"],
                hw_addr=arg.get("mac", ""),
                mtu=arg.get("mtu", 0),
                vrf=arg.get("vrf", 0))
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
                                    eth_src=arg["eth_src"],
                                    phy_port_id=arg.get("phy_port", arg["port"]),
                                    tun_type=arg.get("tun_type", "NOP"),
                                    tun_remote=arg.get("tun_remote", ""),
                                    tun_local=arg.get("tun_local", ""),
            )
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
            grp = pb.MPLSLabelGroup(dst_id=arg["label"],
                                    new_label=arg["out_label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_L3_VPN", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_mpls_tun1(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSLabelGroup(dst_id=arg["label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_TUNNEL1", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    def _mod_group_mpls_swap(self, msg):
        re_id = msg["re_id"]
        cmd = msg["cmd"]
        for arg in msg["args"]:
            grp = pb.MPLSLabelGroup(dst_id=arg["label"],
                                    new_label=arg["out_label"],
                                    ne_id=arg["ne_id"])
            mod = pb.GroupMod(cmd=cmd, g_type="MPLS_SWAP", re_id=re_id, mpls_label=grp)
            evt = fibcevt.EventFIBCGroupMod(mod)

            self.app.send_event_to_observers(evt)


    @staticmethod
    def get_dplist(req, **kwars):
        """
        Get dp_id list.
        """
        _LOG.debug("get_dps %s", req)

        dpids = fibcdbm.dps().keys()
        return Response(content_type='application/json',
                        body=json.dumps(dict(dpids=dpids)))


    def get_port_stats(self, req, dp_id, **kwars):
        """
        Send get stats from dp.
        """
        _LOG.debug("get_stats %s %s", dp_id, req)

        try:
            dpath, mode = fibcdbm.dps().find_by_id(dp_id)
            get_port_stats = ofc.get_port_stats(mode)
            stats = get_port_stats(dpath, self.waiters, None, ofctl)
            _extend_port_stats(stats)
            return Response(content_type='application/json',
                            body=json.dumps(stats))

        except KeyError as ex:
            _LOG.exception(ex)
            return Response(status=404)

        except Exception as ex: # pylint: disable=broad-except
            _LOG.exception(ex)
            return Response(status=505, body=str(ex))


def _extend_port_stats(stats):
    def _snmp_if_oper_status(status):
        return 1 if status else 2

    for dpid, port_stats_list in stats.items():
        dpid = int(dpid)
        for port_stats in port_stats_list:
            try:
                port_no = port_stats["port_no"]
                port_entry = fibcdbm.portmap().find_by_dp(dpid, port_no)
                port_stats["ifName"] = port_entry["name"][1]
                port_stats["ifOperStatus"] = _snmp_if_oper_status(port_entry["dpenter"])

            except KeyError:
                port_stats["ifName"] = ""
                port_stats["ifOperStatus"] = _snmp_if_oper_status(False)
