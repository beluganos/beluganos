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
	fibcnet "fabricflow/fibc/net"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Packet struct {
	Ifindex int
	Data    []byte
	Ffpkt   *fibcnet.FFPacket
}

func NewPacket(ifindex int, data []byte) *Packet {
	return &Packet{
		Ifindex: ifindex,
		Data:    data,
	}
}

func ParseFFPacket(data []byte) (*Packet, error) {
	ffpkt := fibcnet.FFPacket{}
	if err := fibcnet.ParseFFPacket(data, &ffpkt); err != nil {
		return nil, err
	}

	return &Packet{
		Ffpkt: &ffpkt,
	}, nil
}

func StripVlanHeader(eth *layers.Ethernet, vlan *layers.Dot1Q) (*Packet, error) {
	buf := gopacket.NewSerializeBuffer()
	opt := gopacket.SerializeOptions{
		FixLengths: true,
	}

	payload := vlan.LayerPayload()
	d, _ := buf.AppendBytes(len(payload))
	copy(d, payload)

	eth.EthernetType = vlan.Type
	if err := eth.SerializeTo(buf, opt); err != nil {
		return nil, err
	}

	return &Packet{
		Data: buf.Bytes(),
	}, nil
}

func ParsePacket(data []byte, vid uint16) (*Packet, error) {
	pkt := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Lazy)

	var eth *layers.Ethernet

	if layer := pkt.Layer(layers.LayerTypeEthernet); layer == nil {
		return nil, fmt.Errorf("ethernet layer not found.")
	} else {
		eth = layer.(*layers.Ethernet)
	}

	if eth.EthernetType == layers.EthernetType(fibcnet.FFPACKET_ETHTYPE) {
		return ParseFFPacket(data)
	}

	if layer := pkt.Layer(layers.LayerTypeDot1Q); layer != nil {
		if vlan := layer.(*layers.Dot1Q); vlan.VLANIdentifier == vid {
			return StripVlanHeader(eth, vlan)
		}
	}

	return &Packet{
		Data: data,
	}, nil
}
