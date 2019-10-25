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

package gonslib

import (
	fibcapi "fabricflow/fibc/api"

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

const (
	linkmonRegisterName = "linkmon"
)

//
// LinkInfo is opennsl port_info and port_no.
//
type LinkInfo struct {
	opennsl.PortInfo
	Port opennsl.Port
}

//
// NewLinkInfo create new instance of LinkInfo.
//
func NewLinkInfo(port opennsl.Port, info *opennsl.PortInfo) *LinkInfo {
	return &LinkInfo{
		Port:     port,
		PortInfo: *info,
	}
}

//
// PortNo returns port_no.
//
func (l *LinkInfo) PortNo() uint32 {
	return uint32(l.Port)
}

//
// PortState returns port statte.
//
// port is up  : fibcapi.FFPORT_STATE_NONE
// port is down: fibcapi.FFPORT_STATE_LINKDOWN
//
func (l *LinkInfo) PortState() uint32 {
	if l.LinkStatus().IsUp() {
		return fibcapi.FFPORT_STATE_NONE
	}

	return fibcapi.FFPORT_STATE_LINKDOWN
}

//
// LinkmonStart starts link monitor.
//
func (s *Server) LinkmonStart(done <-chan struct{}) <-chan *LinkInfo {
	ch := make(chan *LinkInfo)
	go LinkmonServe(s.Unit(), ch, done)

	return ch
}

//
// LinkmonServe monitor link state and notify.
//
func LinkmonServe(unit int, linkCh chan<- *LinkInfo, done <-chan struct{}) {
	if err := opennsl.LinkscanRegister(unit, linkmonRegisterName, func(unit int, key string, port opennsl.Port, portInfo *opennsl.PortInfo) {
		linkCh <- NewLinkInfo(port, portInfo)
	}); err != nil {
		log.Errorf("LinkscanRegister error. %s", err)
		return
	}

	defer opennsl.LinkscanUnregister(unit, linkmonRegisterName)

	log.Infof("Server: LinkMon: Started.")

	<-done

	log.Infof("Server: LinkMon: Exit")
}
