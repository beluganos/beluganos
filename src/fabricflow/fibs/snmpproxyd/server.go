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

	"github.com/PromonLogicalis/snmp"
	log "github.com/sirupsen/logrus"
)

const (
	SERVER_CHECK_WORKER_INTERVAL = 1 * time.Second
	SERVER_UDP_SEND_TIMEOUT      = 3 * time.Second
	SERVER_UDP_READ_TIMEOUT      = 1 * time.Second
	MIB_NAME_IFINDEX             = "ifIndex"
)

type ProxyServer struct {
	*Tables
	listenAddr *net.UDPAddr
	snmpdAddr  *net.UDPAddr
	snmpComm   string
}

func NewProxyServer(listenAddr, snmpdAddr *net.UDPAddr, ifNotifyCom string) (*ProxyServer, error) {
	return &ProxyServer{
		Tables:     NewTables(),
		listenAddr: listenAddr,
		snmpdAddr:  snmpdAddr,
		snmpComm:   ifNotifyCom,
	}, nil
}

func (s *ProxyServer) Serve() {
	log.Infof("ProxyServer.Serve START")

	ticker := time.NewTicker(SERVER_CHECK_WORKER_INTERVAL)
	for {
		select {
		case <-ticker.C:
			s.WorkerTable().CheckAlive(func(key string, worker *ProxyWorker) {
				worker.Stop()
				log.Debugf("ProxyServer.Serve. worker(%s) stopped.", key)
			})
		}
	}
}

func (s *ProxyServer) Start() {
	go s.Serve()

	lib.StartUDPServer(s.listenAddr, func(buf []byte, clientAddr *net.UDPAddr, clientConn *net.UDPConn) {
		key := clientAddr.String()
		worker, ok := s.WorkerTable().Find(key)
		if !ok {
			worker = NewProxyWorker(clientAddr, clientConn, s)
			if err := worker.Start(); err != nil {
				log.Errorf("ProxyServer.Start worker(%s) start error. %s", key, err)
				return
			}

			log.Debugf("ProxyServer.Start worker(%s) created.", key)
			s.WorkerTable().Put(key, worker)
		}
		worker.Put(buf)
	})
}

func (s *ProxyServer) newSnmpdConn() (conn *net.UDPConn, err error) {
	conn, err = net.DialUDP("udp", nil, s.snmpdAddr)
	return
}

func (s *ProxyServer) convVarToLocal(variable snmp.Variable) snmp.Variable {
	oid := variable.Name

	if len(oid) == 0 {
		log.Warnf("ProxyServer.convVarToLocal invalid oid. %v", oid)
		return variable
	}

	oidmap, ok := s.OidMapTable().MatchByGlobal(oid.String())
	if !ok {
		return variable
	}

	variable.Name, _ = lib.ReplaceOID(oid, oidmap.GlobalOid, oidmap.LocalOid)

	// oidIndex := func() uint {
	//	if len(oid) <= len(oidmap.GlobalOid) {
	//		return 0
	//	}
	//	return oid[len(oid)-1]
	//}()
	//if _, ok := s.IfMapTable().Match(IFMAP_NAME_OIDMAP, oidIndex); ok {
	//	variable.Name, _ = lib.ReplaceOID(oid, oidmap.GlobalOid)
	//}
	//if ifmap, ok := s.IfMapTable().Match(IFMAP_NAME_SHIFT, oidIndex); ok {
	//	oidIndex = oidIndex - ifmap.Min
	//	if oidIndex == 0 {
	//		// remove last. (ex .1.2.3.<index> -> .1.2.3)
	//		variable.Name = variable.Name[:len(variable.Name)-1]
	//	} else {
	//		// replace last to new index. (ex .1.2.3.<index> -> .1.2.3.<index - offset>)
	//		variable.Name[len(variable.Name)-1] = oidIndex
	//	}
	//}

	return variable
}

func (s *ProxyServer) convVarsToLocal(variables []snmp.Variable) []snmp.Variable {
	for index, variable := range variables {
		variables[index] = s.convVarToLocal(variable)
	}

	return variables
}

func (s *ProxyServer) convVarToGlobal(variable snmp.Variable) snmp.Variable {
	oid := variable.Name

	if len(oid) == 0 {
		log.Warnf("ProxyServer.convVarToGlobal invalid oid. %v", oid)
		return variable
	}

	oidmap, ok := s.OidMapTable().MatchByLocal(oid.String())
	if !ok {
		return variable
	}

	variable.Name, _ = lib.ReplaceOID(oid, oidmap.LocalOid, oidmap.GlobalOid)

	//if oidmap, ok := s.OidMapTable().MatchByLocal(oid.String()); ok {
	//	variable.Name, _ = lib.ReplaceOID(oid, oidmap.LocalOid)
	//}
	//if oidmap, ok := s.OidMapTable().MatchByGlobal(oid.String()); ok {
	//	if len(oid) > len(oidmap.GlobalOid) {
	//		ifmap, _ := s.IfMapTable().Get(IFMAP_NAME_SHIFT)
	//		oidIndex := oid[len(oid)-1]
	//		variable.Name[len(variable.Name)-1] = oidIndex + ifmap.Min
	//		if oidmap.Name == MIB_NAME_IFINDEX {
	//			variable.Value = int(oidIndex + ifmap.Min)
	//		}
	//	}
	//}

	return variable
}

func (s *ProxyServer) convVarsToGlobal(variables []snmp.Variable) []snmp.Variable {
	for index, variable := range variables {
		variables[index] = s.convVarToGlobal(variable)
	}

	return variables
}

func (s *ProxyServer) convVarTrap(variable snmp.Variable) snmp.Variable {
	oid := variable.Name

	if len(oid) == 0 {
		log.Warnf("ProxyServer.convVarsIfindex invalid oid. %v", oid)
		return variable
	}

	oidmap, ok := s.OidMapTable().MatchByGlobal(oid.String())
	if !ok {
		return variable
	}
	if len(oid) <= len(oidmap.GlobalOid) {
		return variable // no subfix.
	}

	lastIndex := len(oid) - 1
	ifindex := int(oid[lastIndex])

	if trapmap, ok := s.TrapMapTable().FindByIfindex(ifindex); ok {
		portId := trapmap.PortId
		variable.Name[lastIndex] = uint(portId)

		if oidmap.Name == MIB_NAME_IFINDEX {
			variable.Value = int(portId)
		}
	}
	//else {
	//	ifmap, _ := s.IfMapTable().Get(IFMAP_NAME_SHIFT)
	//	variable.Name[lastIndex] = oid[lastIndex] + ifmap.Min
	//	if oidmap.Name == MIB_NAME_IFINDEX {
	//		variable.Value = int(oid[lastIndex] + ifmap.Min)
	//	}
	//}

	return variable
}

func (s *ProxyServer) convVarsTrap(variables []snmp.Variable) []snmp.Variable {
	for index, variable := range variables {
		variables[index] = s.convVarTrap(variable)
	}

	return variables
}

func (s *ProxyServer) IsPrivateCommunity(c string) bool {
	return c == s.snmpComm
}
