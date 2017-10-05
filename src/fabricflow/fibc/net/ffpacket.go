// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

package fibcnet

import (
	"bytes"
	"encoding/binary"
	"net"
	"syscall"
)

const (
	FFPACKET_ETHTYPE uint16 = 0x0a0a
)

type EthHdr struct {
	EthDst  [6]byte
	EthSrc  [6]byte
	EthType uint16
}

type FFPacket struct {
	EthHdr
	_      uint16
	ReId   [24]byte
	Ifname [24]byte
}

func (p *FFPacket) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, p); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func SockAddr(ifindex int) syscall.Sockaddr {
	return &syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ALL,
		Ifindex:  ifindex,
		Halen:    6,
	}
}

func NewFFPacket(reId string, hwAddr net.HardwareAddr, ifname string) *FFPacket {
	ffpkt := &FFPacket{
		EthHdr: EthHdr{
			EthDst:  [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthType: FFPACKET_ETHTYPE,
		},
	}
	copy(ffpkt.EthSrc[:], hwAddr)
	copy(ffpkt.ReId[:len(ffpkt.ReId)-1], reId)
	copy(ffpkt.Ifname[:len(ffpkt.Ifname)-1], ifname)

	return ffpkt
}
