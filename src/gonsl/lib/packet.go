// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
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

package gonslib

import (
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

const (
	packetRxPri     = 10
	packetRxSize    = 16 * 1024
	packetRxPerChan = 10
	packerRxChans   = 4
	packetRxPps     = 30000
	packetTxPadsize = 4
)

//
// RxInit initialize opennsl rx function.
//
func RxInit(unit int) error {
	pcfg, err := opennsl.PortConfigGet(unit)
	if err != nil {
		log.Errorf("Server: RxInit: PortConfigGet error. %s", err)
		return err
	}

	bmp, _ := pcfg.PBmp(opennsl.PORT_CONFIG_CPU)
	if err := opennsl.VlanDefaultMustGet(unit).PortAdd(unit, bmp, bmp); err != nil {
		log.Errorf("DEFAULT_VLAN.PortAdd. %s", err)
		return err
	}

	log.Infof("Rx init ok.")
	return nil
}

//
// RxStart starts receiving packets.
//
func (s *Server) RxStart(done <-chan struct{}) <-chan *opennsl.Pkt {
	rxCh := make(chan *opennsl.Pkt)
	go RxServe(s.Unit(), rxCh, done)

	return rxCh
}

//
// RxServe serve receiving packets.
//
func RxServe(unit int, rxCh chan<- *opennsl.Pkt, done <-chan struct{}) {

	flg := opennsl.NewRxCallbackFlags(fieldCosDefault)
	err := opennsl.RxRegister(unit, packetRxPri, flg, func(unit int, pkt *opennsl.Pkt) {
		rxCh <- pkt
	})
	if err != nil {
		log.Errorf("Server: RxPacket: RxRegister error. %s", err)
		return
	}

	defer opennsl.RxUnregister(unit, packetRxPri)

	if active := opennsl.RxActive(unit); !active {
		cfg := opennsl.NewRxCfg()
		cfg.SetPktSize(packetRxSize)
		cfg.SetPktsPerChain(packetRxPerChan)
		cfg.SetGlobalPps(packetRxPps)
		cfg.ChanCfg(1).SetChains(packerRxChans)
		// cfg.ChanCfg(1).SetCosBmp(0xffffffff)

		if err := opennsl.RxStart(unit, cfg); err != nil {
			log.Errorf("Server: RxPacket: RxStart error. %s", err)
			return
		}

		defer cfg.Stop(unit)

		log.Infof("Server: RxPacket: activated.")
	}

	log.Infof("Server: RxPacket: Started.")

	<-done

	log.Infof("Server: RxPacket: Exit")
}

//
// FIBCFFPacketOut process FFPacketOUT message from fibcd.
//
func (s *Server) FIBCFFPacketOut(hdr *fibcnet.Header, pktout *fibcapi.FFPacketOut) {
	log.Debugf("Server: PacketOUT: %v %d %d", hdr, pktout.GetPortNo(), len(pktout.Data))

	pkt, err := opennsl.PktAlloc(s.Unit(), len(pktout.Data)+packetTxPadsize, opennsl.PKT_F_NONE)
	if err != nil {
		log.Errorf("Server: PacketOut: PktAlloc error. %s", err)
		return
	}

	defer pkt.Free(s.Unit())

	if err := pkt.Memcpy(0, pktout.GetData()); err != nil {
		log.Errorf("Server: PacketOut: Memcpy error. %s", err)
		return
	}

	pkt.TxPBmp().Clear()
	pkt.TxPBmp().Add(opennsl.Port(pktout.GetPortNo()))

	if err := pkt.Tx(s.Unit()); err != nil {
		log.Errorf("Server: PacketOut: Send error. %s", err)
		return
	}
}
