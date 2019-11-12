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
	"encoding/csv"
	"encoding/hex"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"
)

func parseHexDumpFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseHexDump(f)
}

func parseHexDump(r io.Reader) ([]byte, error) {
	reader := csv.NewReader(r)
	reader.Comma = ' '

	datas := []byte{}
	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			v, err := strconv.ParseUint(record, 16, 8)
			if err != nil {
				return nil, err
			}

			datas = append(datas, uint8(v))
		}
	}

	return datas, nil
}

func hexdumpDebugLog(data []byte) {
	dumper := hex.Dumper(log.StandardLogger().Out)
	defer dumper.Close()

	dumper.Write(data)
}

func newPacketARP() []byte {
	dstHwAddr := net.HardwareAddr([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	srcHwAddr := net.HardwareAddr([]byte{0x00, 0x1a, 0x6b, 0x6c, 0x0c, 0xcc})
	eth := layers.Ethernet{
		DstMAC:       dstHwAddr,
		SrcMAC:       srcHwAddr,
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4, // IPv4
		Operation:         layers.ARPRequest,
		SourceHwAddress:   srcHwAddr,
		SourceProtAddress: net.IP([]byte{10, 0, 0, 1}),
		DstHwAddress:      dstHwAddr,
		DstProtAddress:    net.IP([]byte{10, 0, 0, 2}),
	}
	payload := gopacket.Payload(make([]byte, 86))
	opt := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf, opt,
		&eth,
		&arp,
		payload,
	)

	return buf.Bytes()
}
