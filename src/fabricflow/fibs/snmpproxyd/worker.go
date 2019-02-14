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
	"net"
	"time"

	lib "fabricflow/fibs/fibslib"

	"github.com/PromonLogicalis/asn1"
	"github.com/PromonLogicalis/snmp"
	log "github.com/sirupsen/logrus"
)

const WORKER_MSG_CHAN_SIZE = 16

type ProxyWorker struct {
	*ProxyServer
	ctx        *asn1.Context
	clientAddr *net.UDPAddr
	clientConn *net.UDPConn
	snmpdConn  *net.UDPConn
	msgCh      chan *snmp.Message
	log        *log.Entry
}

func NewProxyWorker(clientAddr *net.UDPAddr, clientConn *net.UDPConn, server *ProxyServer) *ProxyWorker {
	return &ProxyWorker{
		ProxyServer: server,
		ctx:         snmp.Asn1Context(),
		clientAddr:  clientAddr,
		clientConn:  clientConn,
		snmpdConn:   nil,
		msgCh:       make(chan *snmp.Message, WORKER_MSG_CHAN_SIZE),
		log: log.WithFields(log.Fields{
			"client": clientAddr.String(),
		}),
	}
}

func (s *ProxyWorker) Start() error {
	snmpdConn, err := s.newSnmpdConn()
	if err != nil {
		s.log.Errorf("ProxyWorker.SendRecvSnmpd newSnmpConn error. %s", err)
		return err
	}

	s.snmpdConn = snmpdConn

	go s.Serve()

	return nil
}

func (s *ProxyWorker) Stop() {
	s.snmpdConn.Close()
	close(s.msgCh)
}

func (s *ProxyWorker) Put(buf []byte) error {
	msg, err := s.decodeBuffer(buf)
	if err != nil {
		s.log.Errorf("ProxyWorker.Start Decode buffer error. %s", err)
		return err
	}

	s.msgCh <- msg
	return nil
}

func (s *ProxyWorker) Serve() {
	s.log.Debugf("ProxyWorker.Serve START")

	for msg := range s.msgCh {
		s.dispatchMsg(msg)
	}

	s.log.Debugf("ProxyWorker.Serve END")
}

func (s *ProxyWorker) decodeBuffer(buf []byte) (*snmp.Message, error) {
	msg := &snmp.Message{}
	rem, err := s.ctx.Decode(buf, msg)
	if err != nil {
		s.log.Errorf("ProxyWorker Decode error. %s", err)
		return nil, err
	}

	if len(rem) > 0 {
		s.log.Errorf("ProxyWorker Malformed message.")
		return nil, err
	}

	return msg, nil
}

func (s *ProxyWorker) sendMessage(conn *net.UDPConn, msg *snmp.Message, addr *net.UDPAddr) error {
	buf, err := s.ctx.Encode(*msg)
	if err != nil {
		s.log.Errorf("ProxyWorker.sendMessage Encode error. %s", err)
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(SERVER_UDP_SEND_TIMEOUT))

	if addr == nil {
		if _, err := conn.Write(buf); err != nil {
			s.log.Errorf("ProxyWorker.sendMessage Write error. %s", err)
			return err
		}
	} else {
		if _, err := conn.WriteTo(buf, addr); err != nil {
			s.log.Errorf("ProxyWorker.sendMessage WriteTo error. %s", err)
			return err
		}
	}

	return nil
}

func (s *ProxyWorker) recvMessage(conn *net.UDPConn) (*snmp.Message, error) {
	conn.SetReadDeadline(time.Now().Add(SERVER_UDP_READ_TIMEOUT))

	resBuf := make([]byte, lib.UDP_BUFFER_SIZE)
	n, _, err := conn.ReadFromUDP(resBuf)
	if err != nil {
		s.log.Errorf("ProxyWorker.recvMessage Read error. %s", err)
		return nil, err
	}

	s.log.Debugf("ProxyWorker.recvMessage Read %d", n)

	resMsg, err := s.decodeBuffer(resBuf[:n])
	if err != nil {
		s.log.Errorf("ProxyWorker.recvMessage Decode error. %s", err)
		return nil, err
	}

	return resMsg, nil
}

func (s *ProxyWorker) SendRecvSnmpd(msg *snmp.Message) (*snmp.Message, error) {
	if err := s.sendMessage(s.snmpdConn, msg, nil); err != nil {
		s.log.Errorf("ProxyWorker.SendRecvSnmpd sendMessage error. %s", err)
		return nil, err
	}

	return s.recvMessage(s.snmpdConn)
}

func (s *ProxyWorker) SendClient(msg *snmp.Message) error {
	return s.sendMessage(s.clientConn, msg, s.clientAddr)
}

func (s *ProxyWorker) SendTrap(msg *snmp.Message) {
	s.TrapSinkTable().GetAll(func(addr *net.UDPAddr) {
		if err := s.sendMessage(s.clientConn, msg, addr); err != nil {
			s.log.Errorf("ProxyWorker.SendTrap error. %s", err)
		}
	})
}

func (s *ProxyWorker) dispatchMsg(msg *snmp.Message) {
	s.log.Debugf("ProxyWorker.dispatchMsg %v", msg)

	switch pdu := msg.Pdu.(type) {
	case snmp.GetRequestPdu:
		s.processGetRequest(msg, pdu)

	case snmp.GetNextRequestPdu:
		s.processGetNextRequest(msg, pdu)

	case snmp.GetBulkRequestPdu:
		s.processGetBulkRequest(msg, pdu)

	case snmp.SetRequestPdu:
		s.processSetRequest(msg, pdu)

	case snmp.GetResponsePdu:
		s.processGetResponse(msg, pdu)

	case snmp.V1TrapPdu:
		s.processV1Trap(msg, pdu)

	case snmp.V2TrapPdu:
		s.processV2Trap(msg, pdu)

	case snmp.InformRequestPdu:
		s.processInformRequest(msg, pdu)

	default:
		s.log.Errorf("ProxyWorker.dispatchMsg Unsupported PDU. %v", msg.Pdu)
	}
}

func (s *ProxyWorker) processGetRequest(msg *snmp.Message, pdu snmp.GetRequestPdu) {
	s.log.Debugf("ProxyWorker.GetRequest START %v", msg)

	s.log.Debugf("ProxyWorker.GetRequest Request(Global)")
	dumpSnmpPdu((*snmp.Pdu)(&pdu))

	pdu.Variables = s.convVarsToLocal(pdu.Variables)

	s.log.Debugf("ProxyWorker.GetRequest Request(Local)")
	dumpSnmpPdu((*snmp.Pdu)(&pdu))

	msg.Pdu = pdu
	resMsg, err := s.SendRecvSnmpd(msg)
	if err != nil {
		s.log.Errorf("ProxyWorker.GetRequest SendRecvSnmpd error. %s", err)
		return
	}

	resPdu := resMsg.Pdu.(snmp.GetResponsePdu)

	s.log.Debugf("ProxyWorker.GetRequest Response(Local)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resPdu.Variables = s.convVarsToGlobal(resPdu.Variables)

	s.log.Debugf("ProxyWorker.GetRequest Response(Global)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resMsg.Pdu = resPdu
	if err := s.SendClient(resMsg); err != nil {
		s.log.Errorf("ProxyWorker.GetRequest SendClient error. %s", err)
		return
	}

	s.log.Debugf("ProxyWorker.GetRequest END.")
}

func (s *ProxyWorker) processGetNextRequest(msg *snmp.Message, pdu snmp.GetNextRequestPdu) {
	s.log.Debugf("ProxyWorker.GetNextRequest START %v", msg)

	s.log.Debugf("ProxyWorker.GetNextRequest Request(Global)")
	dumpSnmpPdu((*snmp.Pdu)(&pdu))

	pdu.Variables = s.convVarsToLocal(pdu.Variables)

	s.log.Debugf("ProxyWorker.GetNextRequest Request(Local)")
	dumpSnmpPdu((*snmp.Pdu)(&pdu))

	msg.Pdu = pdu
	resMsg, err := s.SendRecvSnmpd(msg)
	if err != nil {
		s.log.Errorf("ProxyWorker.GetNextRequest SendRecvSnmpd error. %s", err)
		return
	}

	resPdu := resMsg.Pdu.(snmp.GetResponsePdu)

	s.log.Debugf("ProxyWorker.GetNextRequest Response(Local)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resPdu.Variables = s.convVarsToGlobal(resPdu.Variables)

	s.log.Debugf("ProxyWorker.GetNextRequest Response(Global)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resMsg.Pdu = resPdu
	if err := s.SendClient(resMsg); err != nil {
		s.log.Errorf("ProxyWorker.GetNextRequest SendClient error. %s", err)
		return
	}

	s.log.Debugf("ProxyWorker.GetNextRequest END.")
}

func (s *ProxyWorker) processGetBulkRequest(msg *snmp.Message, pdu snmp.GetBulkRequestPdu) {
	s.log.Debugf("ProxyWorker.GetBulkRequest START %v", msg)

	s.log.Debugf("ProxyWorker.GetBulkRequest Request(Global)")
	dumpSnmpBulkPdu((*snmp.BulkPdu)(&pdu))

	pdu.Variables = s.convVarsToLocal(pdu.Variables)

	s.log.Debugf("ProxyWorker.GetBulkRequest Request(Local)")
	dumpSnmpBulkPdu((*snmp.BulkPdu)(&pdu))

	msg.Pdu = pdu
	resMsg, err := s.SendRecvSnmpd(msg)
	if err != nil {
		s.log.Errorf("ProxyWorker.GetBulkRequest SendRecvSnmpd error. %s", err)
		return
	}

	resPdu := resMsg.Pdu.(snmp.GetResponsePdu)

	s.log.Debugf("ProxyWorker.GetBulkRequest Response(Local)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resPdu.Variables = s.convVarsToGlobal(resPdu.Variables)

	s.log.Debugf("ProxyWorker.GetBulkRequest Response(Global)")
	dumpSnmpPdu((*snmp.Pdu)(&resPdu))

	resMsg.Pdu = resPdu
	if err := s.SendClient(resMsg); err != nil {
		s.log.Errorf("ProxyWorker.GetBulkRequest SendClient error. %s", err)
		return
	}

	s.log.Debugf("ProxyWorker.GetBulkRequest END")
}

func (s *ProxyWorker) processSetRequest(msg *snmp.Message, pdu snmp.SetRequestPdu) {
	s.log.Debugf("ProxyWorker.SetRequest START %v", msg)
	dumpSnmpPdu((*snmp.Pdu)(&pdu))
	s.log.Debugf("ProxyWorker.SetRequest END")
}

func (s *ProxyWorker) processGetResponse(msg *snmp.Message, pdu snmp.GetResponsePdu) {
	s.log.Debugf("ProxyWorker.GetResponse START %v", msg)
	dumpSnmpPdu((*snmp.Pdu)(&pdu))
	s.log.Debugf("ProxyWorker.GetResponse END")
}

func (s *ProxyWorker) processV1Trap(msg *snmp.Message, pdu snmp.V1TrapPdu) {
	s.log.Debugf("ProxyWorker.V1Trap START %v", msg)
	dumpSnmpTrapV1Pdu(&pdu)
	s.log.Debugf("ProxyWorker.V1Trap END")
}

func (s *ProxyWorker) processV2Trap(msg *snmp.Message, pdu snmp.V2TrapPdu) {
	s.log.Debugf("ProxyWorker.V2Trap START %v", msg)

	if s.IsPrivateCommunity(msg.Community) {
		s.updateIfindexMapTable(pdu)
	} else {
		s.log.Debugf("ProxyWorker.V2Trap (Local)")
		dumpSnmpPdu((*snmp.Pdu)(&pdu))

		pdu.Variables = s.convVarsTrap(pdu.Variables)

		s.log.Debugf("ProxyWorker.V2Trap Request(Global)")
		dumpSnmpPdu((*snmp.Pdu)(&pdu))

		msg.Pdu = pdu
		s.SendTrap(msg)
	}

	s.log.Debugf("ProxyWorker.V2Trap END")
}

func (s *ProxyWorker) updateIfindexMapTable(pdu snmp.V2TrapPdu) {
	for _, variable := range pdu.Variables {
		oid := variable.Name
		if len(oid) == 0 {
			log.Errorf("ProxyWorker.V2Trap iface invalid oid. %v", variable)
			continue
		}

		ifname, ok := variable.Value.(string)
		if !ok {
			s.log.Errorf("ProxyWorker.V2Trap iface invalid var. %v", variable)
			continue
		}

		s.TrapMapTable().UpdateByIfname(ifname, func(e *TrapMapEntry) {
			e.Ifindex = int(oid[len(oid)-1])
			s.log.Infof("ProxyWorker.V2Trap iface '%s' ifindex:%d", ifname, e.Ifindex)
		})
	}
}

func (s *ProxyWorker) processInformRequest(msg *snmp.Message, pdu snmp.InformRequestPdu) {
	s.log.Debugf("ProxyWorker.InformRequest START %v", msg)
	dumpSnmpPdu((*snmp.Pdu)(&pdu))
	s.log.Debugf("ProxyWorker.InformRequest END")
}
