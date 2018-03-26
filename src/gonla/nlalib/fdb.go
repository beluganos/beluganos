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

package nlalib

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	FDB_SYSFS_CLASS_NET = "/sys/class/net"
)

type FdbEntry struct {
	MacAddr     [6]byte
	PortNo      uint8
	Local       uint8
	AgeingTimer uint32
	PortHi      uint8
	Pad0        uint8
	Unused      uint16
}

func (e *FdbEntry) HardwareAddr() net.HardwareAddr {
	return net.HardwareAddr(e.MacAddr[:])
}

func (e *FdbEntry) Port() uint16 {
	return (uint16(e.PortHi) << 8) + uint16(e.PortNo)
}

func (e *FdbEntry) IsLocal() bool {
	return e.Local != 0
}

func (e *FdbEntry) String() string {
	return fmt.Sprintf("%s %d local:%t age:%d", e.HardwareAddr(), e.Port(), e.IsLocal(), e.AgeingTimer)
}

func FdbSysfsPath(ifname string) string {
	return fmt.Sprintf("%s/%s/brforward", FDB_SYSFS_CLASS_NET, ifname)
}

func ReadFdb(ifname string) ([]*FdbEntry, error) {
	f, err := os.Open(FdbSysfsPath(ifname))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries := []*FdbEntry{}
	for {
		entry := &FdbEntry{}
		err := binary.Read(f, binary.LittleEndian, entry)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
