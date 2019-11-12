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
	"encoding/hex"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

func dumpPktDesc(pkt *opennsl.Pkt) {
	log.Debugf("rx: cos:%d port:%d vid:%d len:%d tot:%d",
		pkt.Cos(), pkt.RxPort(), pkt.VID(), pkt.PktLen(), pkt.TotLen())
}

func dumpPktDetail(pkt *opennsl.Pkt) {
	log.Debugf("pkt  : %p len:%d tot:%d", pkt, pkt.PktLen(), pkt.TotLen())
	log.Debugf("unit : %d", pkt.Unit())
	log.Debugf("flags: %d", pkt.Flags())
	log.Debugf("cos  : %d", pkt.Cos())
	log.Debugf("vid  : %d", pkt.VID())
	log.Debugf("port : src:%d dst:%d", pkt.SrcPort(), pkt.DstPort())
	log.Debugf("rx   : port    : %d", pkt.RxPort())
	log.Debugf("rx   : untagged: %d", pkt.RxUntagged())
	log.Debugf("rx   : matched : %d", pkt.RxMatched())
	log.Debugf("rx   : reasons : %d", pkt.RxReasons())
	log.Debugf("blk  : #%d", pkt.BlkCount())
}

func dumpPktData(pkt *opennsl.Pkt) {
	total := int(pkt.PktLen())
	if total > 256 {
		total = 256
	}

	for index, blk := range pkt.Blks() {
		dumpLen := func() int {
			if blkLen := blk.Len(); blkLen < total {
				return blkLen
			}
			return total
		}()

		total -= dumpLen

		log.Debugf("blk[%d] len=%d", index, dumpLen)
		b := blk.Data()
		log.Debugf("\n%s", hex.Dump(b[:dumpLen]))
	}
}

//
// dumpRxPkt output packet.
//
func (s *Server) dumpRxPkt(pkt *opennsl.Pkt) {
	if s.LogConfig().RxSilent {
		return
	}

	if s.LogConfig().RxDetail {
		dumpPktDetail(pkt)
	} else {
		dumpPktDesc(pkt)
	}

	if s.LogConfig().RxDump {
		dumpPktData(pkt)
	}
}

//
// dumpTxPkt output packet.
//
func (s *Server) dumpTxPkt(pkt *opennsl.Pkt) {
	if s.LogConfig().TxSilent {
		return
	}

	if s.LogConfig().TxDetail {
		dumpPktDetail(pkt)
	} else {
		dumpPktDesc(pkt)
	}

	if s.LogConfig().TxDump {
		dumpPktData(pkt)
	}
}
