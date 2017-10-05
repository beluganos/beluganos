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

package nlaapi

import (
	"gonla/nladbm"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"net"
)

func (n *Node) GetIP() net.IP {
	return net.IP(n.Ip)
}

func (n *Node) ToNetlink() *nlalink.Node {
	return &nlalink.Node{
		Ip: n.Ip,
	}
}

func (n *Node) ToNative() *nlamsg.Node {
	return &nlamsg.Node{
		Node: n.ToNetlink(),
		NId:  uint8(n.NId),
		Ch:   nil,
	}
}

func NewNodeFromNative(n *nlamsg.Node) *Node {
	return &Node{
		Ip:  n.Ip,
		NId: uint32(n.NId),
	}
}

//
// Node (Key)
//
func (k *NodeKey) ToNative() *nladbm.NodeKey {
	return nladbm.NewNodeKey(uint8(k.NId))
}

func NewNodeKeyFromNative(n *nladbm.NodeKey) *NodeKey {
	return &NodeKey{
		NId: uint32(n.NId),
	}
}
