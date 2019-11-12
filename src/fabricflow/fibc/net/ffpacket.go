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
	"fmt"
	"net"
	"syscall"
)

const (
	FFPACKET_ETHTYPE uint16 = 0x0a0a

	FFPACKET_DATA_SIZE = 6 + 6 + 2 + 2 + 24 + 24
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

func ffpkBytesToString(b []byte) string {
	if index := bytes.IndexByte(b, 0); index >= 0 {
		return string(b[0:index])
	}

	return string(b)
}

func (d *FFPacket) GetReId() string {
	return ffpkBytesToString(d.ReId[:])
}

func (d *FFPacket) GetIfname() string {
	return ffpkBytesToString(d.Ifname[:])
}

func ParseFFPacket(data []byte, ffpkt *FFPacket) error {
	if l := len(data); l < FFPACKET_DATA_SIZE {
		return fmt.Errorf("Invalid length. len=%d", l)
	}

	copy(ffpkt.EthDst[:], data[0:6])
	copy(ffpkt.EthSrc[:], data[6:12])
	ffpkt.EthType = uint16(data[13])<<8 + uint16(data[12])
	copy(ffpkt.ReId[:], data[16:40])
	copy(ffpkt.Ifname[:], data[40:64])

	return nil
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
