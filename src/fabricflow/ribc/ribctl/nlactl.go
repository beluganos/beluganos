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

package ribctl

import (
	fibcapi "fabricflow/fibc/api"
	"gonla/nlaapi"
	"gonla/nlalib"
	"gonla/nlamsg"
	"io"
	"net"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type NLAController struct {
	addr   string
	connCh chan net.IP
	recvCh chan *nlamsg.NetlinkMessageUnion
	client nlaapi.NLAApiClient
	log    *log.Entry
}

func NewNLAController(addr string) *NLAController {
	return &NLAController{
		addr:   addr,
		connCh: make(chan net.IP),
		recvCh: make(chan *nlamsg.NetlinkMessageUnion),
		client: nil,
		log:    log.WithFields(log.Fields{"module": "NLAController"}),
	}
}

func (n *NLAController) Conn() <-chan net.IP {
	return n.connCh
}

func (n *NLAController) Recv() <-chan *nlamsg.NetlinkMessageUnion {
	return n.recvCh
}

func (n *NLAController) Start() error {
	n.log.Debugf("Start:")

	ch := make(chan *nlalib.ConnInfo)
	conn, err := nlalib.NewClientConn(n.addr, ch)
	if err != nil {
		n.log.Errorf("Start: connect error. %s", err)
		close(ch)
		return err
	}
	n.client = nlaapi.NewNLAApiClient(conn)

	go func() {
		for {
			ci := <-ch
			go n.Monitor(ci)
		}
	}()

	return nil
}

func (n *NLAController) Monitor(ci *nlalib.ConnInfo) {
	n.connCh <- ci.LocalAddr
	defer func() {
		n.connCh <- nil
	}()

	stream, err := n.client.MonNetlink(context.Background(), &nlaapi.MonNetlinkRequest{})
	if err != nil {
		n.log.Errorf("Monitor: error. %s", err)
		return
	}

	n.log.Infof("Monitor: START")

	for {
		nlmsg, err := stream.Recv()
		if err != nil {
			n.log.Infof("Monitor: EXIT. %s", err)
			break
		}

		n.log.Debugf("Monitor: recv %v", nlmsg)
		n.recvCh <- nlmsg.ToNative()
	}
}

func (n *NLAController) GetLink(nid uint8, index int) (*nlamsg.Link, error) {
	key := &nlaapi.LinkKey{
		NId:   uint32(nid),
		Index: int32(index),
	}
	link, err := n.client.GetLink(context.Background(), key)
	if err != nil {
		return nil, err
	}

	return link.ToNative(), nil
}

func (n *NLAController) GetLink_GroupMod(cmd fibcapi.GroupMod_Cmd, nid uint8, index int) (*nlamsg.Link, error) {
	link, err := n.GetLink(nid, index)
	if cmd == fibcapi.GroupMod_DELETE {
		err = nil
	}

	return link, err
}

func (n *NLAController) GetLinks(nid uint8, f func(*nlamsg.Link) error) error {
	stream, err := n.client.GetLinks(context.Background(), nlaapi.NewGetLinksRequest(nid))
	if err != nil {
		return err
	}

	for {
		link, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := f(link.ToNative()); err != nil {
			return err
		}
	}
}

func (n *NLAController) GetIptun(nid uint8, remote net.IP) (*nlamsg.Iptun, error) {
	key := &nlaapi.IptunKey{
		NId:    uint32(nid),
		Remote: remote,
	}
	iptun, err := n.client.GetIptun(context.Background(), key)
	if err != nil {
		return nil, err
	}

	return iptun.ToNative(), nil
}

func (n *NLAController) GetAddrs(nid uint8, f func(*nlamsg.Addr) error) error {
	stream, err := n.client.GetAddrs(context.Background(), nlaapi.NewGetAddrsRequest(nid))
	if err != nil {
		return nil
	}

	for {
		addr, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := f(addr.ToNative()); err != nil {
			return nil
		}
	}
}

func (n *NLAController) GetNeigh(nid uint8, addr net.IP) (*nlamsg.Neigh, error) {
	key := &nlaapi.NeighKey{
		NId:  uint32(nid),
		Addr: addr.String(),
	}
	neigh, err := n.client.GetNeigh(context.Background(), key)
	if err != nil {
		return nil, err
	}
	return neigh.ToNative(), nil
}

func (n *NLAController) GetNeigh_FlowMod(cmd fibcapi.FlowMod_Cmd, nid uint8, addr net.IP) (*nlamsg.Neigh, error) {
	neigh, err := n.GetNeigh(nid, addr)
	if cmd == fibcapi.FlowMod_DELETE || cmd == fibcapi.FlowMod_DELETE_STRICT {
		err = nil
	}

	return neigh, err
}

func (n *NLAController) GetNeigh_GroupMod(cmd fibcapi.GroupMod_Cmd, nid uint8, addr net.IP) (*nlamsg.Neigh, error) {
	neigh, err := n.GetNeigh(nid, addr)
	if cmd == fibcapi.GroupMod_DELETE {
		err = nil
	}

	return neigh, err
}

func (n *NLAController) GetNeighs(nid uint8, f func(*nlamsg.Neigh) error) error {
	stream, err := n.client.GetNeighs(context.Background(), nlaapi.NewGetNeighsRequest(nid))
	if err != nil {
		return nil
	}

	for {
		neigh, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := f(neigh.ToNative()); err != nil {
			return nil
		}
	}
}

func (n *NLAController) GetRoutes(nid uint8, f func(*nlamsg.Route) error) error {
	stream, err := n.client.GetRoutes(context.Background(), nlaapi.NewGetRoutesRequest(nid))
	if err != nil {
		return nil
	}

	for {
		route, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := f(route.ToNative()); err != nil {
			return nil
		}
	}
}

func (n *NLAController) GetBridgeVlanInfos(nid uint8, ifindex int, f func(*nlamsg.BridgeVlanInfo)) error {
	stream, err := n.client.GetBridgeVlanInfos(context.Background(), nlaapi.NewGetBridgeVlanInfosRequest(nid))
	if err != nil {
		return err
	}

	for {
		brvlan, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		v := brvlan.ToNative()
		if nid == v.NId && ifindex == v.Index {
			f(v)
		}
	}
}

func (n *NLAController) ModLinkStatus(nid uint8, ifname string, operState string) error {
	link := nlaapi.NewDeviceLink(nid, 0)
	attr := link.GetDevice().GetLinkAttrs()
	attr.Name = ifname
	attr.OperState = nlaapi.ParseLinkOperState(operState)

	req := nlaapi.NewNetlinkMessageUnion(nid, syscall.RTM_SETLINK, link)
	_, err := n.client.ModNetlink(context.Background(), req)

	return err
}

func (n *NLAController) ModFdb(nid uint8, ifindex int, hwaddr net.HardwareAddr, vid uint16, mtype uint16) error {
	neigh := &nlaapi.Neigh{
		NId:          uint32(nid),
		Ip:           []byte{},
		LinkIndex:    int32(ifindex),
		HardwareAddr: hwaddr,
		VlanId:       int32(vid),
	}

	req := nlaapi.NewNetlinkMessageUnion(nid, mtype, neigh)
	_, err := n.client.ModNetlink(context.Background(), req)

	return err
}
