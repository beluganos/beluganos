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
	"fmt"
	"gonla/nlaapi"
	"gonla/nladbm"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//
// NLACoreApiServer
//
type NLACoreApiServer struct {
	addr   string
	NlMsgs chan<- *nlamsg.NetlinkMessage
	log    *log.Entry
}

func NewNLACoreApiServer(addr string) *NLACoreApiServer {
	return &NLACoreApiServer{
		addr:   addr,
		NlMsgs: nil,
		log:    NewLogger("NLACoreApiServer"),
	}
}

func (n *NLACoreApiServer) Start(ch chan<- *nlamsg.NetlinkMessage) error {
	listen, err := net.Listen("tcp", n.addr)
	if err != nil {
		return err
	}

	n.NlMsgs = ch

	s := grpc.NewServer()
	nlaapi.RegisterNLACoreApiServer(s, n)
	go s.Serve(listen)

	n.log.Infof("Start:")
	return nil
}

func (n *NLACoreApiServer) SendNetlinkMessage(ctxt context.Context, req *nlaapi.NetlinkMessage) (*nlaapi.NetlinkMessageReply, error) {
	n.NlMsgs <- req.ToNative()
	return &nlaapi.NetlinkMessageReply{}, nil
}

func (n *NLACoreApiServer) MonNetlinkMessage(req *nlaapi.Node, stream nlaapi.NLACoreApi_MonNetlinkMessageServer) error {
	n.log.Infof("Monitor: START %v", req)

	node := req.ToNative()
	nid := node.NId

	sendNode := func(t uint16) {
		nlmsg, _ := nlamsg.NodeSerialize(node, t)
		n.NlMsgs <- nlmsg
	}

	if old := nladbm.Nodes().Insert(node); old != nil {
		return fmt.Errorf("NLACoreApiServer: Node already exists. %v", old)
	}

	node.Open()
	defer node.Close()

	sendNode(nlalink.RTM_NEWNODE)

	if done := stream.Context().Done(); done != nil {
		go func() {
			<-done
			n.log.Infof("Monitor:  EXIT nid:%d", nid)

			nladbm.Nodes().Delete(nladbm.NodeToKey(node))
			sendNode(nlalink.RTM_DELNODE)

			node.Send(nil)
		}()
	}

	for m := range node.Recv() {
		if m == nil {
			break
		}

		n.log.Debugf("Monitor; send to slave. nid:%d %v", nid, m)
		nlmsg := nlaapi.NewNetlinkMessageUnionFromNative(m)
		if err := stream.Send(nlmsg); err != nil {
			n.log.Errorf("Monitor: stream error. nid:%d %s", nid, err)
		}
	}

	return nil
}
