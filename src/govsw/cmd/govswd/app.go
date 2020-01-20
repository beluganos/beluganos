// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/hex"
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"
	"fmt"
	"govsw/pkgs/govsw"
	"io"

	log "github.com/sirupsen/logrus"
)

type MyApp struct {
	db     *govsw.DB
	cfg    *govsw.ConfigServer
	syncCh chan string
	client fibcapi.FIBCVsApiClient

	log *log.Entry
}

func NewMyApp(db *govsw.DB, cfg *govsw.ConfigServer) *MyApp {
	return &MyApp{
		db:     db,
		cfg:    cfg,
		syncCh: make(chan string),
		log:    log.WithField("module", "app"),
	}
}

func (a *MyApp) SyncCh() <-chan string {
	return a.syncCh
}

func (a *MyApp) VswAPISync(name string) {
	go func(ifname string) {
		a.log.Infof("db sync '%s'", ifname)
		a.syncCh <- ifname
	}(name)
}

func (a *MyApp) VswAPISaveConfig() error {
	ifnames := []string{}
	patterns := []string{}
	a.db.Ifname().Range(func(kind string, name string) {
		switch kind {
		case "name":
			ifnames = append(ifnames, name)
		case "pattern":
			patterns = append(patterns, name)
		}
	})

	a.log.Debugf("VswAPISaveConfig: ifnames :%v", ifnames)
	a.log.Debugf("VswAPISaveConfig: patterns:%v", patterns)

	a.cfg.SetIfNames(ifnames)
	a.cfg.SetIfPatterns(patterns)

	if err := a.cfg.Write(); err != nil {
		a.log.Errorf("VswAPISave: Write error. %s", err)
		return err
	}

	return nil
}

func (a *MyApp) ConfigChanged(cfg *govsw.DpConfig) {
	a.db.SetDpID(cfg.DpId)
	a.db.Ifname().Update(cfg)
	a.log.Infof("db updated")
}

func (a *MyApp) PacketIn(pkt *govsw.Packet) {
	if pkt.Ffpkt != nil {
		reId := pkt.Ffpkt.GetReId()
		ifname := pkt.Ffpkt.GetIfname()

		a.log.Tracef("PacketIn: FFPacket(%d %s %s)", pkt.Ifindex, reId, ifname)

		if err := a.sendFFPacket(pkt.Ifindex, reId, ifname); err != nil {
			a.log.Errorf("PacketIn: send ffpacket error, %s", err)
		}

	} else {
		if log.IsLevelEnabled(log.TraceLevel) {
			dumpSize := len(pkt.Data)
			if dumpSize > 64 {
				dumpSize = 64
			}
			a.log.Tracef("PacketIn: ifindex=%d size=%d", pkt.Ifindex, len(pkt.Data))
			a.log.Tracef("PacketIn:\n%s", hex.Dump(pkt.Data[:dumpSize]))
		}

		if err := a.sendPacketIn(pkt.Ifindex, pkt.Data); err != nil {
			a.log.Errorf("PacketIn: send packet error. %s", err)
		}
	}
}

func (a *MyApp) FIBCClient(client fibcapi.FIBCVsApiClient) {
	a.client = client
}

func (a *MyApp) FIBCConnect() {
	if err := a.sendHello(); err != nil {
		a.log.Errorf("FIBClient: hello error. %s", err)
		return
	}

	go a.monitor()
}

func (a *MyApp) sendHello() error {
	hello := fibcapi.FFHello{
		DpId:   a.db.DpID(),
		DpType: fibcapi.FFHello_FFVS,
	}

	_, err := a.client.SendHello(context.Background(), &hello)
	return err
}

func (a *MyApp) sendPacketIn(ifindex int, data []byte) error {
	pktin := fibcapi.FFPacketIn{
		DpId:   a.db.DpID(),
		PortNo: uint32(ifindex),
		Data:   data,
	}

	_, err := a.client.SendPacketIn(context.Background(), &pktin)
	return err
}

func (a *MyApp) sendFFPacket(ifindex int, reId, ifname string) error {
	ffpkt := fibcapi.FFPacket{
		DpId:   a.db.DpID(),
		PortNo: uint32(ifindex),
		ReId:   reId,
		Ifname: ifname,
	}

	_, err := a.client.SendFFPacket(context.Background(), &ffpkt)
	return err
}

func (a *MyApp) sendOAMReply(reply *fibcapi.OAM_Reply, xid uint32) error {
	msg := fibcapi.OAMReply{
		Xid:   xid,
		Reply: reply,
	}

	_, err := a.client.SendOAMReply(context.Background(), &msg)
	return err
}

func (a *MyApp) monitor() {
	monreq := fibcapi.VsMonitorRequest{
		VsId:   a.db.DpID(),
		DpType: fibcapi.FFHello_FFVS,
	}
	stream, err := a.client.Monitor(context.Background(), &monreq)
	if err != nil {
		a.log.Errorf("monitor: Monitor error. %s", err)
		return
	}

	a.log.Debugf("monitor: started.")

FOR_LOOP:
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			a.log.Infof("monitor: exit.")
			break FOR_LOOP
		}
		if err != nil {
			a.log.Errorf("monitor: exit. recv error. %s", err)
			break FOR_LOOP
		}
		fibcapi.DispatchVsMonitorReply(m, a)
	}
}

func (a *MyApp) FIBCFFPacketOut(hdr *fibcnet.Header, pktout *fibcapi.FFPacketOut) {
	a.log.Tracef("monitor: packet out dpid:%d port:%d #%d",
		pktout.DpId,
		pktout.PortNo,
		len(pktout.Data),
	)

	if err := a.packetOut(pktout); err != nil {
		a.log.Errorf("monitor: packet out error. %s", err)
	}
}

func (a *MyApp) packetOut(pktout *fibcapi.FFPacketOut) error {
	return a.db.Link().Get(int(pktout.PortNo), func(link *govsw.Link) error {
		link.WriteData(pktout.Data)
		return nil
	})
}

func (a *MyApp) FIBCFFPortMod(hdr *fibcnet.Header, pm *fibcapi.FFPortMod) {
	a.log.Debugf("monitor: port mod dpid:%d port:%d %s",
		pm.DpId,
		pm.PortNo,
		pm.Status,
	)

	if err := a.portMod(pm); err != nil {
		a.log.Errorf("monitor: port mod error. %s", err)
	}
}

func (a *MyApp) portMod(mod *fibcapi.FFPortMod) error {
	return a.db.Link().Get(int(mod.PortNo), func(link *govsw.Link) error {
		switch mod.Status {
		case fibcapi.PortStatus_UP:
			return link.SetOperStatus(true)

		case fibcapi.PortStatus_DOWN:
			return link.SetOperStatus(false)

		default:
			return fmt.Errorf("Invaid status. %s", mod.Status)
		}
	})
}

func (a *MyApp) FIBCOAMAuditRouteCntRequest(hdr *fibcnet.Header, oam *fibcapi.OAM_Request, audit *fibcapi.OAM_AuditRouteCntRequest) error {
	msg := fibcapi.OAM_Reply{
		DpId:    a.db.DpID(),
		OamType: oam.OamType,
		Body: &fibcapi.OAM_Reply_AuditRouteCnt{
			AuditRouteCnt: &fibcapi.OAM_AuditRouteCntReply{
				Count: 0,
			},
		},
	}
	go a.sendOAMReply(&msg, hdr.Xid)
	return nil
}
