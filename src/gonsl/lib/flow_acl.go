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
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

//
// FIBCPolicyACLFlowMod process FlowMod(ACL Policy)
//
func (s *Server) FIBCPolicyACLFlowMod(hdr *fibcnet.Header, mod *fibcapi.FlowMod, flow *fibcapi.PolicyACLFlow) {
	log.Debugf("Server: FlowMod(ACL): %v %v", hdr, mod)

	switch {
	case len(flow.Match.IpDst) != 0:
		log.Debugf("Server: FlowMod(ACL): ip_dst. %s", flow)

		_, dstIP, err := net.ParseCIDR(flow.Match.IpDst)
		if err != nil {
			ip := net.ParseIP(flow.Match.IpDst)
			if ip == nil {
				log.Errorf("Server: FlowMod(ACL): Invalid IP. %s", flow.Match.IpDst)
				return
			}

			bits := EtherTypeToLen(uint16(flow.Match.EthType))
			mask := net.CIDRMask(bits, bits)
			dstIP = &net.IPNet{
				IP:   ip,
				Mask: mask,
			}
		}

		switch mod.Cmd {
		case fibcapi.FlowMod_ADD:
			err := func() error {
				switch flow.Match.EthType {
				case unix.ETH_P_IP:
					return s.Fields().DstIPv4.AddEntry(dstIP)
				case unix.ETH_P_IPV6:
					return s.Fields().DstIPv6.AddEntry(dstIP)
				default:
					return fmt.Errorf("Invalid ether type. %04x", flow.Match.EthType)
				}
			}()
			if err != nil {
				log.Errorf("Server: FlowMod(ACL): AddEntry error. %s %s", dstIP, err)
			}

		case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
			switch flow.Match.EthType {
			case unix.ETH_P_IP:
				s.Fields().DstIPv4.DeleteEntry(dstIP)
			case unix.ETH_P_IPV6:
				s.Fields().DstIPv6.DeleteEntry(dstIP)
			default:
				log.Warnf("Server: FlowMod(ACL): DeleteEntry error. %s", dstIP)
			}

		default:
			log.Warnf("Server: FlowMod(ACL): Invalid cmd. %s", mod.Cmd)
		}

	case flow.Match.EthType != 0:
		log.Debugf("Server: FlowMod(ACL): eth_type. %s", flow)

		switch mod.Cmd {
		case fibcapi.FlowMod_ADD:
			if err := s.Fields().EthType.AddEntry(uint16(flow.Match.EthType)); err != nil {
				log.Errorf("Server: FlowMod(ACL): AddEntry error. %d %s", flow.Match.EthType, err)
			}

		case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
			s.Fields().EthType.DeleteEntry(uint16(flow.Match.EthType))

		default:
			log.Warnf("Server: FlowMod(ACL): Invalid cmd. %s", mod.Cmd)
		}

	case len(flow.Match.EthDst) > 0:
		log.Debugf("Server:  FlowMod(ACL): eth_dst. %s", flow)

		dstMAC, err := net.ParseMAC(flow.Match.EthDst)
		if err != nil {
			log.Errorf("Server: FlowMod(ACL): Invalid MAC. %s", flow.Match.EthDst)
			return
		}

		switch mod.Cmd {
		case fibcapi.FlowMod_ADD:
			if err := s.Fields().EthDst.AddEntry(dstMAC, fibcapi.HardwareAddrExactMask); err != nil {
				log.Errorf("Server:  FlowMod(ACL): AddEntry error. %s %s", dstMAC, err)
			}

		case fibcapi.FlowMod_DELETE, fibcapi.FlowMod_DELETE_STRICT:
			s.Fields().EthDst.DeleteEntry(dstMAC, fibcapi.HardwareAddrExactMask)

		default:
			log.Warnf("Server: FlowMod(ACL): Invalid cmd. %s", mod.Cmd)
		}

	default:
		log.Warnf("Server: FlowMod(ACL): Ignored. %s", flow)
	}
}
