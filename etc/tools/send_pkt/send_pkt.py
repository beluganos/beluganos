#! /usr/bin/env python
# -*- coding: utf-8 -*-

import yaml
import copy
import traceback
from scapy.all import (TCP, UDP, IP, IPv6, Ether,
                       Dot1Q, sendp, ARP, ICMP, load_contrib)
load_contrib("ospf")

def is_proto(proto, s):
    if isinstance(proto, str):
        return proto == s

    if isinstance(proto, list) and len(proto) > 0:
        return proto[0] == s

    return False

def get_ipver(s):
    if s is not None and ":" in s:
        return 5
    return 4

class TxPayload(object):
    def __init__(self, cfg):
        self.cfg = cfg

    def pkt(self):
        args = copy.copy(self.cfg)
        proto = args.pop("proto")

        if is_proto(proto, "arp"):
            return ARP(**args)

        dst = args.pop("dst", None)
        src = args.pop("src", None)

        ipver = get_ipver(dst)
        if ipver == 4:
            ip = IP(dst=dst, src=src)
        else:
            ip = IPv6(dst=dst, src=src)

        if is_proto(proto, "tcp"):
            return ip / TCP(**args)

        elif is_proto(proto, "udp"):
            return ip / UDP(**args)

        elif is_proto(proto, "icmp"):
            return ip / ICMP(**args)

        elif is_proto(proto, "ospf"):
            msg = proto[1]
            if msg == "hello":
                if ipver == 4:
                    return ip / OSPF_Hdr() / OSPF_Hello(**args)
                else:
                    return ip / OSPFv3_Hdr() / OSPFv3_Hello(**args)

            elif msg == "dbdesc":
                if ipver == 4:
                    return ip / OSPF_Hdr() / OSPF_DBDesc(**args)
                else:
                    return ip / OSPFv3_Hdr() / OSPFv3_DBDesc(**args)


class TxPacket(object):
    def __init__(self, cfg):
        self.cfg = cfg

    def pkt(self):
        pkt = Ether(dst=self.cfg["dst"], src=self.cfg.get("src", None))
        vid = self.cfg.get("vid", None)
        if vid is not None:
            pkt = pkt / Dot1Q(vlan=vid)

        return pkt


class TxConfig(object):
    def __init__(self, cfg):
        self.cfg = cfg

    def sends(self):
        return self.cfg.get("send", list())

    def payloads(self):
        payloads = self.cfg.get("payload", dict())
        return {name:TxPayload(payload) for name, payload in payloads.items()}

    def ifaces(self):
        return self.cfg.get("iface", dict())

    def packets(self):
        packets = self.cfg.get("packet", dict())
        return {name:TxPacket(packet) for name, packet in packets.items()}


class TxRunner(object):
    def __init__(self, cfg, dry_run = False):
        self.cfg = cfg
        self.count = 0
        self.dry_run = dry_run

    def get_cnt(self):
        self.count += 1
        return self.count

    def make_data(self, data):
        return "cnt={0} {1}".format(self.get_cnt(), data)

    def run(self):
        for send in self.cfg.sends():
            try:
                iface = self.cfg.ifaces().get(send["iface"], "lo")
                data = self.make_data("if={0}".format(iface))
                payload = self.cfg.payloads().get(send["payload"])
                pkt = self.cfg.packets().get(send["packet"])
                pkt = pkt.pkt() / payload.pkt() / data
                pkt.show2()

                if not self.dry_run:
                    sendp(pkt, iface=iface)

            except:
                traceback.print_exc()


def read_config(path):
    import yaml
    with open(path) as f:
        return yaml.load(f, Loader=yaml.Loader)

def _args():
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config")
    parser.add_argument("--dry-run", action="store_true")
    return parser.parse_args()

def _main():
    args = _args()

    cfg = read_config(args.config)
    tx = TxRunner(TxConfig(cfg["tx"]), args.dry_run)
    tx.run()

if __name__ == "__main__":
    _main()
