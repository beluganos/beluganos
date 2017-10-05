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
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gonla/nlaapi"
	"gonla/nladbm"
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"google.golang.org/grpc"
	"net"
)

//
// NLACoreApiServer
//
type NLACoreApiServer struct {
	addr   string
	NlMsgs chan<- *nlamsg.NetlinkMessage
}

func NewNLACoreApiServer(addr string) *NLACoreApiServer {
	return &NLACoreApiServer{
		addr:   addr,
		NlMsgs: nil,
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

	log.Infof("NLACoreApiServer: START")
	return nil
}

func (n *NLACoreApiServer) SendNetlinkMessage(ctxt context.Context, req *nlaapi.NetlinkMessage) (*nlaapi.NetlinkMessageReply, error) {
	n.NlMsgs <- req.ToNative()
	return &nlaapi.NetlinkMessageReply{}, nil
}

func (n *NLACoreApiServer) MonNetlinkMessage(req *nlaapi.Node, stream nlaapi.NLACoreApi_MonNetlinkMessageServer) error {
	log.Infof("NLACoreApiServer: Monitor START. %v", req)

	node := req.ToNative()
	nid := node.NId

	sendNode := func(t uint16) {
		nlmsg, _ := nlamsg.NodeSerialize(node, t)
		n.NlMsgs <- nlmsg
	}

	if old := nladbm.Nodes().Insert(node); old != nil {
		return fmt.Errorf("NLACoreApiServer: Node already exists. %v", old)
	}
	defer nladbm.Nodes().Delete(nladbm.NodeToKey(node))

	sendNode(nlalink.RTM_NEWNODE)
	defer sendNode(nlalink.RTM_DELNODE)

	done := stream.Context().Done()
	for {
		select {
		case <-done:
			log.Infof("NLACoreApiServer: Monitor EXIT. nid:%d", nid)
			return nil

		case m := <-node.Recv():
			log.Debugf("NLACoreApiServer: Send to slave. nid:%d %v", nid, m)
			nlmsg := nlaapi.NewNetlinkMessageUnionFromNative(m)
			if err := stream.Send(nlmsg); err != nil {
				log.Errorf("NLACoreApiServer: Stream.Send error. nid:%d %s", nid, err)
				return err
			}
		}
	}

	// log.Infof("Api Server; Monitor EXIT %v", req)
	// return nil
}
