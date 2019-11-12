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

package govsw

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/google/gopacket/pcap"
	"github.com/vishvananda/netlink"

	log "github.com/sirupsen/logrus"
)

const (
	LINK_BLOCK_SIZE      = 1024 * 1024
	LINK_PCAPBUF_SIZE    = 16 * 1024
	LINK_WRITE_CHAN_SIZE = 64
	LINK_VLAN_ID_DEFAULT = 1
)

const (
	linkStatsPacketRead        = "pkt/read"
	linkStatsPacketReadEnQ     = "pkt/read/enq"
	linkStatsPacketReadErr     = "pkt/read/err"
	linkStatsPacketWrite       = "pkt/write"
	linkStatsPacketWriteEnQ    = "pkt/write/enq"
	linkStatsPacketWriteErr    = "pkt/write/err"
	linkStatsOperStatusUP      = "oper/up"
	linkStatsOperStatusUPErr   = "oper/up/err"
	linkStatsOperStatusDown    = "oper/down"
	linkStatsOperStatusDownErr = "oper/down/err"
)

var linkStatsNames = []string{
	linkStatsPacketRead,
	linkStatsPacketReadEnQ,
	linkStatsPacketReadErr,
	linkStatsPacketWrite,
	linkStatsPacketWriteEnQ,
	linkStatsPacketWriteErr,
	linkStatsOperStatusUP,
	linkStatsOperStatusUPErr,
	linkStatsOperStatusDown,
	linkStatsOperStatusDownErr,
}

func NewLinkStats() *StatsGroup {
	stats := NewStatsGroup()
	stats.RegisterList(linkStatsNames)
	return stats
}

type Link struct {
	name  string
	index int

	socket     *pcap.Handle
	writeCh    chan []byte
	writeStrip uint16
	readStrip  uint16

	done  chan struct{}
	mutex sync.Mutex
	stats *StatsGroup

	log *log.Entry
}

func NewLink(link netlink.Link) *Link {
	return &Link{
		name:  link.Attrs().Name,
		index: link.Attrs().Index,

		writeCh: make(chan []byte, LINK_WRITE_CHAN_SIZE),
		stats:   NewLinkStats(),

		log: log.WithField("module", fmt.Sprintf("link/%d", link.Attrs().Index)),
	}
}

func (l *Link) String() string {
	return fmt.Sprintf("ifname:'%s' ifindex:%d", l.name, l.index)
}

func (l *Link) Name() string {
	return l.name
}

func (l *Link) Index() int {
	return l.index
}

func (l *Link) Stats() *StatsGroup {
	return l.stats
}

func (l *Link) SetStripWPkt(v uint16) {
	l.writeStrip = v
}

func (l *Link) SetStripRPkt(v uint16) {
	l.readStrip = v
}

//func (l *Link) newTPacket() (*afpacket.TPacket, error) {
//	return afpacket.NewTPacket(
//		afpacket.OptInterface(l.name),
//		afpacket.OptBlockSize(LINK_BLOCK_SIZE),
//	)
//}

func (l *Link) newPcapHandle() (*pcap.Handle, error) {
	return pcap.OpenLive(
		l.name,
		LINK_PCAPBUF_SIZE,
		true,              // promisc mode
		pcap.BlockForever, // timeout
	)
}

func (l *Link) readPacket(sock *pcap.Handle, ch chan<- *Packet) {
	l.log.Debugf("readPacket: Started")

	defer l.Stop()

FOR_LOOP:
	for {
		l.stats.Inc(linkStatsPacketRead)

		data, _, err := sock.ReadPacketData()
		if err != nil {
			l.stats.Inc(linkStatsPacketReadErr)

			l.log.Errorf("readPacket: socket read error. %s", err)
			break FOR_LOOP
		}

		if log.IsLevelEnabled(log.TraceLevel) {
			log.Tracef("readPacket: size:%d", len(data))
			log.Tracef("readPacket: \n%s", hex.Dump(data))
		}

		if l.readStrip > 0 {
			dataLen := uint16(len(data))
			if dataLen > l.readStrip {
				data = data[:dataLen-l.readStrip]
			}
		}

		pkt, err := ParsePacket(data, LINK_VLAN_ID_DEFAULT)
		if err != nil {
			l.stats.Inc(linkStatsPacketReadErr)

			l.log.Errorf("readPacket: Parse packet error. %s", err)
			continue
		}

		pkt.Ifindex = l.index
		ch <- pkt

		l.stats.Inc(linkStatsPacketReadEnQ)
	}

	l.log.Debugf("readPacket: Exit")
}

func (l *Link) writePacket(sock *pcap.Handle, done chan struct{}) {
	l.log.Debugf("writePacket: Started")

FOR_LOOP:
	for {
		select {
		case data := <-l.writeCh:
			l.stats.Inc(linkStatsPacketWrite)

			if log.IsLevelEnabled(log.TraceLevel) {
				log.Tracef("writePacket: size:%d strip:%d", len(data), l.writeStrip)
				log.Tracef("writePacket: \n%s", hex.Dump(data))
			}

			if l.writeStrip > 0 {
				if dataLen := uint16(len(data)); dataLen > l.writeStrip {
					data = data[:(dataLen - l.writeStrip)]
				}
			}

			// strip vlan(vid=1) header.
			pkt, err := ParsePacket(data, LINK_VLAN_ID_DEFAULT)
			if err != nil {
				l.stats.Inc(linkStatsPacketWriteErr)

				l.log.Errorf("writeacket: Parse packet error. %s", err)
				continue
			}

			if err := sock.WritePacketData(pkt.Data); err != nil {
				l.stats.Inc(linkStatsPacketWriteErr)

				l.log.Errorf("writePacket: Write packet error. %s", err)
				// ignore error and continue.
			}

		case <-done:
			break FOR_LOOP
		}
	}

	l.log.Debugf("writePacket: Exit")
}

func (l *Link) waitStop(done <-chan struct{}) {
	<-done
	l.log.Infof("waitStop:")
	l.Stop()
}

func (l *Link) Start(ch chan<- *Packet) error {
	l.log.Infof("Start:")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.done != nil {
		l.log.Warnf("Start: Already started.")
		return nil
	}

	sock, err := l.newPcapHandle()
	if err != nil {
		l.log.Errorf("Start: new socket error. %s", err)
		return err
	}

	l.socket = sock
	l.done = make(chan struct{})

	go l.readPacket(sock, ch)
	go l.writePacket(sock, l.done)
	go l.waitStop(l.done)

	l.log.Infof("Start: end")
	return nil
}

func (l *Link) Stop() {
	l.log.Infof("Stop:")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.done != nil {
		close(l.done)
		l.done = nil
		l.log.Debugf("Stop: channel closed.")
	}

	if l.socket != nil {
		l.socket.Close()
		l.socket = nil
		l.log.Debugf("Stop: socket closed.")
	}

	l.log.Debugf("Stop: end")
}

func (l *Link) Destroy() {
	l.log.Infof("Destroy:")
	l.Stop()
	close(l.writeCh)
}

func (l *Link) SetOperStatus(up bool) error {
	link := netlink.Device{}
	link.Index = l.index

	if up {
		l.stats.Inc(linkStatsOperStatusUP)
		if err := netlink.LinkSetUp(&link); err != nil {
			l.stats.Inc(linkStatsOperStatusUPErr)
			return err
		}
	} else {
		l.stats.Inc(linkStatsOperStatusDown)
		if err := netlink.LinkSetDown(&link); err != nil {
			l.stats.Inc(linkStatsOperStatusDownErr)
			return err
		}
	}

	return nil
}

func (l *Link) WriteData(data []byte) {
	l.writeCh <- data

	l.stats.Inc(linkStatsPacketWriteEnQ)
}
