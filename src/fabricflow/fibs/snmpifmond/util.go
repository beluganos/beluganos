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

package main

import (
	"fmt"
	"net"

	lib "fabricflow/fibs/fibslib"

	"github.com/PromonLogicalis/asn1"
	"github.com/PromonLogicalis/snmp"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type LinkStatus int

const (
	LinkStatusNone LinkStatus = iota
	LinkStatusUp
	LinkStatusDown
)

var linkStatus_names = map[LinkStatus]string{
	LinkStatusNone: "None",
	LinkStatusUp:   "Up",
	LinkStatusDown: "Down",
}

var linkStatus_values = map[string]LinkStatus{
	"None": LinkStatusNone,
	"Up":   LinkStatusUp,
	"Down": LinkStatusDown,
}

func (v LinkStatus) String() string {
	if s, ok := linkStatus_names[v]; ok {
		return s
	}
	return fmt.Sprintf("LinkStatus(%d)", v)
}

func ParseLinkStatus(s string) (LinkStatus, error) {
	if v, ok := linkStatus_values[s]; ok {
		return v, nil
	}
	return LinkStatusNone, fmt.Errorf("Invalid LinkStatus. %s", s)
}

func linkUpdateStatus(msg *netlink.LinkUpdate) LinkStatus {
	switch msg.Header.Type {
	case unix.RTM_NEWLINK:
		return linkStatus(msg.Link)

	case unix.RTM_DELLINK:
		return LinkStatusDown

	default:
		return LinkStatusNone
	}
}

func linkStatus(link netlink.Link) LinkStatus {
	flags := link.Attrs().Flags
	if v := flags & net.FlagUp; v != 0 {
		return LinkStatusUp
	}
	return LinkStatusDown
}

type SnmpMessageBuilder struct {
	version   int
	community string
	ctx       *asn1.Context
}

func NewSnmpMessageBuilder(version int, community string) *SnmpMessageBuilder {
	return &SnmpMessageBuilder{
		version:   version,
		community: community,
		ctx:       snmp.Asn1Context(),
	}
}

func (b *SnmpMessageBuilder) Encode(msg *snmp.Message) ([]byte, error) {
	return b.ctx.Encode(*msg)
}

func (b *SnmpMessageBuilder) NewLinksTrap(oid []uint, links []netlink.Link) (*snmp.Message, error) {
	vars := make([]snmp.Variable, len(links))
	for index, link := range links {
		vars[index] = snmp.Variable{
			Name:  append(lib.CloneOID(oid), uint(link.Attrs().Index)),
			Value: link.Attrs().Name,
		}
	}
	pdu := snmp.V2TrapPdu{
		Variables: vars,
	}

	msg := &snmp.Message{
		Version:   b.version,
		Community: b.community,
		Pdu:       pdu,
	}

	return msg, nil
}
