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

package nlasvc

import (
	"gonla/nlalib"
	"gonla/nlamsg"
	"syscall"
)

func DatasToNetlnkMessages(nid uint8, mtype uint16, datas [][]byte) []*nlamsg.NetlinkMessage {
	msgs := []*nlamsg.NetlinkMessage{}
	for _, data := range datas {
		msg := nlalib.NewNetlinkMessage(mtype, data)
		msgs = append(msgs, nlamsg.NewNetlinkMessage(msg, nid, nlamsg.SRC_KNL))
	}

	return msgs
}

func GetNetlinkMessageLinks(nid uint8) ([]*nlamsg.NetlinkMessage, error) {
	datas, err := nlalib.GetNetlinkLinks()
	if err != nil {
		return nil, err
	}

	return DatasToNetlnkMessages(nid, syscall.RTM_NEWLINK, datas), nil
}

func GetNetlinkMessageAddrs(nid uint8) ([]*nlamsg.NetlinkMessage, error) {
	datas, err := nlalib.GetNetlinkAddrs()
	if err != nil {
		return nil, err
	}

	return DatasToNetlnkMessages(nid, syscall.RTM_NEWADDR, datas), nil
}

func GetNetlinkMessageNeighs(nid uint8) ([]*nlamsg.NetlinkMessage, error) {
	datas, err := nlalib.GetNetlinkNeighs()
	if err != nil {
		return nil, err
	}

	return DatasToNetlnkMessages(nid, syscall.RTM_NEWNEIGH, datas), nil
}

func GetNetlinkMessageRoutes(nid uint8) ([]*nlamsg.NetlinkMessage, error) {
	datas, err := nlalib.GetNetlinkRoutes()
	if err != nil {
		return nil, err
	}

	return DatasToNetlnkMessages(nid, syscall.RTM_NEWROUTE, datas), nil
}

func SubscribeNetlinkResources(ch chan<- *nlamsg.NetlinkMessage, nid uint8) error {
	funcs := []func(uint8) ([]*nlamsg.NetlinkMessage, error){
		GetNetlinkMessageLinks,
		GetNetlinkMessageAddrs,
		GetNetlinkMessageNeighs,
		GetNetlinkMessageRoutes,
	}
	for _, f := range funcs {
		msgs, err := f(nid)
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			ch <- msg
		}
	}

	return nil
}
