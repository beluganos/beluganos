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
	fibcnet "fabricflow/fibc/net"

	log "github.com/sirupsen/logrus"
)

//
// FIBCGroupMod process any GroupMod.
//
func (s *Server) FIBCGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod) {
	// s.log.Debugf("GroupMod: %v %v", hdr, mod)
	fibcapi.DispatchGroupMod(hdr, mod, s)
}

//
// FIBCMPLSInterfaceGroupMod process GroupMod(MPLS Interface)
//
func (s *Server) FIBCMPLSInterfaceGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSInterfaceGroup) {
	s.log.Debugf("GroupMod(MPLS-IF): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}

//
// FIBCMPLSLabelL2VpnGroupMod process GroupMod(MPLS Label(L2 VPN))
//
func (s *Server) FIBCMPLSLabelL2VpnGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSLabelGroup) {
	s.log.Debugf("GroupMod(MPLS-L2VPN): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}

//
// FIBCMPLSLabelL3VpnGroupMod process GroupMod(MPLS Label(L3 VPN))
//
func (s *Server) FIBCMPLSLabelL3VpnGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSLabelGroup) {
	s.log.Debugf("GroupMod(MPLS-L3VPN): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}

//
// FIBCMPLSLabelTun1GroupMod process GroupMod(MPLS Label(Tunnel1))
//
func (s *Server) FIBCMPLSLabelTun1GroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSLabelGroup) {
	s.log.Debugf("GroupMod(MPLS-Tun1): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}

//
// FIBCMPLSLabelTun2GroupMod process GroupMod(MPLS Label(Tunnel2))
//
func (s *Server) FIBCMPLSLabelTun2GroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSLabelGroup) {
	s.log.Debugf("GroupMod(MPLS-Tun2): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}

//
// FIBCMPLSLabelSwapGroupMod process GroupMod(MPLS Label(Swap))
//
func (s *Server) FIBCMPLSLabelSwapGroupMod(hdr *fibcnet.Header, mod *fibcapi.GroupMod, group *fibcapi.MPLSLabelGroup) {
	s.log.Debugf("GroupMod(MPLS-Swap): %v", hdr)
	fibcapi.LogGroupMod(s.log, log.DebugLevel, mod)
}
