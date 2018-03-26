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

package nlamsg

import (
	"fmt"
	"gonla/nlalib"
	"gonla/nlamsg/nlalink"
)

const NODE_ID_ALL uint8 = 255

func CopyNode(src *nlalink.Node) *nlalink.Node {
	dst := *src
	return &dst
}

//
// Node
//
type Node struct {
	*nlalink.Node
	NId uint8
	Ch  chan *NetlinkMessageUnion
}

func (n *Node) Copy() *Node {
	return &Node{
		Node: CopyNode(n.Node),
		NId:  n.NId,
		Ch:   n.Ch,
	}
}

func (n *Node) Open() {
	n.Ch = make(chan *NetlinkMessageUnion)
}

func (n *Node) Close() {
	if n.Ch != nil {
		close(n.Ch)
		n.Ch = nil
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("{%s} NId: %d", n.IP(), n.NId)
}

func (n *Node) Recv() <-chan *NetlinkMessageUnion {
	return n.Ch
}

func (n *Node) Send(msg *NetlinkMessageUnion) error {
	if msg.NId == n.NId || msg.NId == NODE_ID_ALL {
		n.Ch <- msg
	}
	return nil
}

func NewNode(node *nlalink.Node, nid uint8) *Node {
	return &Node{
		Node: node,
		NId:  nid,
		Ch:   nil,
	}
}

func NodeDeserialize(nlmsg *NetlinkMessage) (*Node, error) {
	node, err := nlalink.NodeDeserialize(nlmsg.Data)
	if err != nil {
		return nil, err
	}

	return NewNode(node, nlmsg.NId), nil
}

func NodeSerialize(node *Node, msgType uint16) (*NetlinkMessage, error) {
	data, err := nlalink.NodeSerialize(node.Node)
	if err != nil {
		return nil, err
	}
	return NewNetlinkMessage(nlalib.NewNetlinkMessage(msgType, data), node.NId, SRC_NOP), nil
}
