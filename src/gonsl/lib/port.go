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

	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

//
// PortBmpGet returns port map instance.
//
func PortBmpGet(unit int) (*opennsl.PBmp, error) {
	portcfg, err := opennsl.PortConfigGet(unit)
	if err != nil {
		return nil, err
	}

	return portcfg.PBmp(opennsl.PORT_CONFIG_E)
}

//
// PortDefaultConfig setup ports.
//
func PortDefaultConfig(unit int) error {
	pbmp, err := PortBmpGet(unit)
	if err != nil {
		return err
	}

	portInfo := opennsl.NewPortInfo()
	portInfo.SetSpeed(0)
	portInfo.SetDuplex(opennsl.PORT_DUPLEX_FULL)
	portInfo.SetPauseRX(opennsl.PORT_ABILITY_PAUSE_RX)
	portInfo.SetPauseTX(opennsl.PORT_ABILITY_PAUSE_TX)
	portInfo.SetLinkscan(opennsl.LINKSCAN_MODE_SW)
	portInfo.SetAutoNeg(false)
	portInfo.SetEnable(true)
	portInfo.SetActionMask(opennsl.NewPortAttr(
		opennsl.PORT_ATTR_AUTONEG_MASK,
		opennsl.PORT_ATTR_DUPLEX_MASK,
		opennsl.PORT_ATTR_PAUSE_TX_MASK,
		opennsl.PORT_ATTR_PAUSE_RX_MASK,
		opennsl.PORT_ATTR_LINKSCAN_MASK,
		opennsl.PORT_ATTR_ENABLE_MASK,
		opennsl.PORT_ATTR_SPEED_MASK,
	))
	stg := opennsl.Stg(1)

	return pbmp.Each(func(port opennsl.Port) error {
		if err := stg.StpSet(unit, port, opennsl.STG_STP_FORWARD); err != nil {
			log.Errorf("StgSet error. %d %s", port, err)
			return err
		}

		if err := port.SelectiveSet(unit, portInfo); err != nil {
			log.Errorf("SelectiveSet error. %d %s", port, err)
			return err
		}

		log.Infof("PortDefaultConfig %d ok.", port)
		return nil
	})
}

//
// PortDefaultVlanConfig registers port to default vlan.
//
func PortDefaultVlanConfig(unit int) error {
	pbmp, err := PortBmpGet(unit)
	if err != nil {
		log.Errorf("PortBmpGet error. %s", err)
		return err
	}

	if err := opennsl.VlanDefaultMustGet(unit).PortAdd(unit, pbmp, pbmp); err != nil {
		log.Errorf("PortAdd error. %s", err)
		return err
	}

	return pbmp.Each(func(port opennsl.Port) error {
		if err := port.UntaggedVlanSet(unit, opennsl.VlanDefaultMustGet(unit)); err != nil {
			log.Errorf("UntaggedVlanSet error. %d %s", port, err)
			return err
		}

		log.Infof("UntaggedVlanSet ok. %d", port)
		return nil
	})
}

//
// PortInfoGet returns opennsl port info.
//
func PortInfoGet(unit int, port opennsl.Port) (*opennsl.PortInfo, error) {
	info := opennsl.NewPortInfo()
	if err := info.PortSelectiveGet(unit, port); err != nil {
		return nil, err
	}

	return info, nil
}

func (s *Server) fibcFFMultipartPortList(portNo uint32) ([]opennsl.Port, error) {
	if portNo != 0xffffffff {
		return []opennsl.Port{opennsl.Port(portNo)}, nil
	}

	pbmp, err := PortBmpGet(s.Unit())
	if err != nil {
		s.log.Errorf("Multipart(Port): PortBmpGet error. %s", err)
		return nil, err
	}

	return pbmp.PortList(), nil
}

func (s *Server) fibcFFMultipartPortRequestGet(hdr *fibcnet.Header, req *fibcapi.FFMultipart_PortRequest) {
	ports, err := s.fibcFFMultipartPortList(req.PortNo)
	if err != nil {
		return
	}

	statsList, err := NewPortStats(req.Names).GetAll(s.Unit(), ports)
	if err != nil {
		s.log.Errorf("Multipart(Port): PortStatsGetAll error. %s", err)
		return
	}

	ffstats := make([]*fibcapi.FFPortStats, len(statsList))
	for index, stats := range statsList {
		ffstats[index] = fibcapi.NewFFPortStats(uint32(ports[index]), stats)
	}

	reply := fibcapi.NewFFMultipart_Reply_Port(s.DpID(), ffstats)
	if err := s.client.MultipartReply(reply, hdr.Xid); err != nil {
		s.log.Errorf("Multipart(Port): Write error. %s", err)
		return
	}
}

func (s *Server) fibcFFMultipartPortRequestReset(hdr *fibcnet.Header, req *fibcapi.FFMultipart_PortRequest) {
	ports, err := s.fibcFFMultipartPortList(req.PortNo)
	if err != nil {
		return
	}

	for _, portID := range ports {
		s.log.Debugf("Multipart(Port): StatClear port=%d", portID)

		if err := opennsl.StatClear(s.Unit(), portID); err != nil {
			s.log.Errorf("Multipart(Port): StatClear error. port=%d %s", portID, err)
		}
	}
}

//
// FIBCFFMultipartPortRequest process multipart-request(port stats) from fibcd.
//
func (s *Server) FIBCFFMultipartPortRequest(hdr *fibcnet.Header, mp *fibcapi.FFMultipart_Request, req *fibcapi.FFMultipart_PortRequest) {
	s.log.Debugf("Multipart(Port): %v", hdr)
	fibcapi.LogFFMultipartRequest(s.log, log.DebugLevel, mp, hdr.Xid)

	switch req.Cmd {
	case fibcapi.FFPortStats_GET:
		s.fibcFFMultipartPortRequestGet(hdr, req)

	case fibcapi.FFPortStats_RESET:
		s.fibcFFMultipartPortRequestReset(hdr, req)

	default:
		s.log.Errorf("Multipart(Port): unsupported cmd. %s", req.Cmd)
	}

	s.log.Debugf("Multipart(Port): end.")
}

//
// FIBCFFMultipartPortDescRequest process multipart-request(port desc) from fibcd.
//
func (s *Server) FIBCFFMultipartPortDescRequest(hdr *fibcnet.Header, mp *fibcapi.FFMultipart_Request, pd *fibcapi.FFMultipart_PortDescRequest) {
	s.log.Debugf("Multipart(PortDesc): %v", hdr)
	fibcapi.LogFFMultipartRequest(s.log, log.DebugLevel, mp, hdr.Xid)

	pbmp, err := PortBmpGet(s.Unit())
	if err != nil {
		s.log.Errorf("Multipart(PortDesc): PortBmpGet error. %s", err)
		return
	}

	ffports := []*fibcapi.FFPort{}
	err = pbmp.Each(func(port opennsl.Port) error {
		linkStatus, _ := port.LinkStatusGet(s.Unit())
		curSpeed, _ := port.SpeedGet(s.Unit())
		maxSpeed, _ := port.SpeedMax(s.Unit())
		hwaddr, _ := port.PauseAddrGet(s.Unit())

		ffport := fibcapi.NewFFPort(uint32(port))
		ffport.CurrSpeed = uint32(curSpeed)
		ffport.MaxSpeed = uint32(maxSpeed)
		ffport.HwAddr = hwaddr.String()
		ffport.State = func() uint32 {
			if linkStatus == 0 {
				return 0x1 // OFPPS_LINK_DOWN
			}
			return 0
		}()

		ffports = append(ffports, ffport)
		return nil
	})
	if err != nil {
		return
	}

	reply := fibcapi.NewMultipart_Reply_PortDesc(s.DpID(), ffports, pd.Internal)
	if err := s.client.MultipartReply(reply, hdr.Xid); err != nil {
		s.log.Errorf("Multipart(PortDesc): Write error. %s", err)
		return
	}

	s.log.Debugf("Multipart(PortDesc): end.")
}

func (s *Server) FIBCFFPortMod(hdr *fibcnet.Header, mod *fibcapi.FFPortMod) {
	s.log.Debugf("PortMod: %v %v", hdr, mod)

	portId, portType := fibcapi.ParseDPPortId(mod.PortNo)
	if portType.IsVirtual() {
		s.log.Debugf("PortMod: skip virtual device. port=%x", mod.PortNo)
		return
	}

	port := opennsl.Port(portId)

	switch mod.Status {
	case fibcapi.PortStatus_UP:
		if err := port.EnableSet(s.Unit(), opennsl.PORT_ENABLE_TRUE); err != nil {
			s.log.Errorf("PortMod: PortEnableSet(TRUE) error. port %d. %s", port, err)
		}

	case fibcapi.PortStatus_DOWN:
		if err := port.EnableSet(s.Unit(), opennsl.PORT_ENABLE_FALSE); err != nil {
			s.log.Errorf("PortMod: PortEnableSet(FALSE) error. port %d. %s", port, err)
		}

	default:
	}
}
