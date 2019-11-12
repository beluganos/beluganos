// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package fibcsrv

import (
	"errors"
	fibcapi "fabricflow/fibc/api"
	"fabricflow/fibc/pkgs/fibcdbm"
	"fmt"
)

var ENoEffect = errors.New("ENoEffect")

//
// Convertor is interface to convert message.
//
type Convertor interface {
	ConvertIDVMtoDP(string) (uint64, error)
	ConvertIDDPtoVM(uint64) (string, error)
	ConvertPortVMtoDP(string, uint32) (uint64, uint32, error)
	ConvertPortDPtoVM(uint64, uint32) (string, uint32, error)
	ConvertPortDPtoVS(uint64, uint32) (uint64, uint32, error)
	ConvertPortVStoDP(uint64, uint32) (uint64, uint32, error)
}

//
// ConvertIDVMtoDP converts reID to dpID
//
func (c *DBCtl) ConvertIDVMtoDP(reID string) (uint64, error) {
	if dpID, ok := c.IDMap().SelectByReID(reID); ok {
		return dpID, nil
	}

	return 0, fmt.Errorf("reid not found. %s", reID)
}

//
// ConvertIDDPtoVM converts dpID to reID.
//
func (c *DBCtl) ConvertIDDPtoVM(dpID uint64) (string, error) {
	if reID, ok := c.IDMap().SelectByDpID(dpID); ok {
		return reID, nil
	}

	return "", fmt.Errorf("dpid not found. %d", dpID)
}

//
// ConvertPortVMtoDP convers portID vm to dp.
//
func (c *DBCtl) ConvertPortVMtoDP(reID string, vmPort uint32) (uint64, uint32, error) {
	var (
		dpID   uint64
		dpPort uint32
	)

	if vmPort == 0 {
		dpID, err := c.ConvertIDVMtoDP(reID)
		if err != nil {
			return 0, 0, err
		}

		return dpID, 0, nil
	}

	ok := c.PortMap().SelectByVM(reID, vmPort, func(e *fibcdbm.PortEntry) {
		dpID = e.DPPort.DpID
		dpPort = e.DPPort.PortID
	})

	if !ok {
		return 0, 0, fmt.Errorf("vm port not found. reid:%s port:%d", reID, vmPort)
	}

	return dpID, dpPort, nil
}

//
// ConvertPortDPtoVM convers portID dp to vm
//
func (c *DBCtl) ConvertPortDPtoVM(dpID uint64, dpPort uint32) (string, uint32, error) {
	var (
		reID   string
		vmPort uint32
	)

	if dpPort == 0 {
		reID, err := c.ConvertIDDPtoVM(dpID)
		if err != nil {
			return "", 0, err
		}

		return reID, 0, nil
	}

	ok := c.PortMap().SelectByDP(dpID, dpPort, func(e *fibcdbm.PortEntry) {
		reID = e.VMPort.ReID
		vmPort = e.VMPort.PortID
	})

	if !ok {
		return "", 0, fmt.Errorf("dp port not found. dpid:%d port:%d", dpID, dpPort)
	}

	return reID, vmPort, nil
}

//
// ConvertPortDPtoVS convers portID dp to vs
//
func (c *DBCtl) ConvertPortDPtoVS(dpID uint64, dpPort uint32) (uint64, uint32, error) {
	var (
		vsID   uint64
		vsPort uint32
	)

	ok := c.PortMap().SelectByDP(dpID, dpPort, func(e *fibcdbm.PortEntry) {
		vsID = e.VSPort.DpID
		vsPort = e.VSPort.PortID
	})

	if !ok {
		return 0, 0, fmt.Errorf("dp port not found. dpid:%d port:%d", dpID, dpPort)
	}

	return vsID, vsPort, nil
}

//
// ConvertPortVStoDP convers portID  vs to dp
//
func (c *DBCtl) ConvertPortVStoDP(vsID uint64, vsPort uint32) (uint64, uint32, error) {
	var (
		dpID   uint64
		dpPort uint32
	)

	ok := c.PortMap().SelectByVS(vsID, vsPort, func(e *fibcdbm.PortEntry) {
		dpID = e.DPPort.DpID
		dpPort = e.DPPort.PortID
	})

	if !ok {
		return 0, 0, fmt.Errorf("vs port not found. vsid:%d port:%d", vsID, vsPort)
	}

	return dpID, dpPort, nil
}

//
// ConvertVLANFlow converts vlan flow.
//
func (c *DBCtl) ConvertVLANFlow(reID string, flow *fibcapi.VLANFlow) error {
	m := flow.GetMatch()
	_, dpPort, err := c.ConvertPortVMtoDP(reID, m.InPort)
	if err != nil {
		c.log.Debugf("ConvertVLANFlow: %s", err)
		return err
	}

	m.InPort = dpPort

	return nil
}

//
// ConvertTermMacFlow converts term mac flow.
//
func (c *DBCtl) ConvertTermMacFlow(reID string, flow *fibcapi.TerminationMacFlow) error {
	m := flow.GetMatch()
	_, dpPort, err := c.ConvertPortVMtoDP(reID, m.InPort)
	if err != nil {
		c.log.Debugf("ConvertTermMacFlow: %s", err)
		return err
	}

	m.InPort = dpPort

	return nil
}

//
// ConvertBridgingFlow converts bridge flow.
//
func (c *DBCtl) ConvertBridgingFlow(reID string, flow *fibcapi.BridgingFlow) error {
	a := flow.GetAction()
	if a.Name == fibcapi.BridgingFlow_Action_OUTPUT && a.Value != 0 {
		_, dpPort, err := c.ConvertPortVMtoDP(reID, a.Value)
		if err != nil {
			c.log.Debugf("ConvertBridgingFlow: %s", err)
			return err
		}

		a.Value = dpPort
	}

	return nil
}

//
// ConvertPolicyACLFlow converts policy acl flow.
//
func (c *DBCtl) ConvertPolicyACLFlow(reID string, flow *fibcapi.PolicyACLFlow) error {
	m := flow.GetMatch()
	if m.InPort == 0 {
		c.log.Debugf("ConvertPolicyACLFlow: %s", ENoEffect)
		return ENoEffect
	}

	_, dpPort, err := c.ConvertPortVMtoDP(reID, m.InPort)
	if err != nil {
		c.log.Debugf("ConvertPolicyACLFlow: %s", err)
		return err
	}

	m.InPort = dpPort

	return nil
}

//
// ConvertFlowMod converts flow mod.
//
func (c *DBCtl) ConvertFlowMod(mod *fibcapi.FlowMod) error {
	reID := mod.ReId

	switch e := mod.Entry.(type) {
	case *fibcapi.FlowMod_Vlan:
		return c.ConvertVLANFlow(reID, e.Vlan)

	case *fibcapi.FlowMod_TermMac:
		return c.ConvertTermMacFlow(reID, e.TermMac)

	case *fibcapi.FlowMod_Mpls1:
		return nil

	case *fibcapi.FlowMod_Unicast:
		return nil

	case *fibcapi.FlowMod_Bridging:
		return c.ConvertBridgingFlow(reID, e.Bridging)

	case *fibcapi.FlowMod_Acl:
		return c.ConvertPolicyACLFlow(reID, e.Acl)

	default:
		return fmt.Errorf("Invalid flow mod. %s %v", reID, e)
	}
}

//
// ConvertL2IfaceGroup converts l2 interface group.
//
func (c *DBCtl) ConvertL2IfaceGroup(reID string, g *fibcapi.L2InterfaceGroup) error {
	_, dpPort, err := c.ConvertPortVMtoDP(reID, g.PortId)
	if err != nil {
		c.log.Debugf("ConvertL2IfaceGroup/port: %s", err)
		return err
	}

	_, dpMaster, err := c.ConvertPortVMtoDP(reID, g.Master)
	if err != nil {
		c.log.Debugf("ConvertL2IfaceGroup/master: %s", err)
		return err
	}

	g.PortId = dpPort
	g.Master = dpMaster

	return nil
}

//
// ConvertL3UnicastGroup converts l3 unicast group.
//
func (c *DBCtl) ConvertL3UnicastGroup(reID string, g *fibcapi.L3UnicastGroup) error {
	_, dpPort, err := c.ConvertPortVMtoDP(reID, g.PortId)
	if err != nil {
		c.log.Debugf("ConvertL3UnicastGroup/port: %s", err)
		return err
	}

	_, dpPhyPort, err := c.ConvertPortVMtoDP(reID, g.PhyPortId)
	if err != nil {
		c.log.Debugf("ConvertL3UnicastGroup/phy: %s", err)
		return err
	}

	g.PortId = dpPort
	g.PhyPortId = dpPhyPort

	return nil
}

//
// ConvertMplsIfaceGroup converts mpls interface group.
//
func (c *DBCtl) ConvertMplsIfaceGroup(reID string, g *fibcapi.MPLSInterfaceGroup) error {
	_, dpPort, err := c.ConvertPortVMtoDP(reID, g.PortId)
	if err != nil {
		c.log.Debugf("ConvertMplsIfaceGroup: %s", err)
		return err
	}

	g.PortId = dpPort

	return nil
}

//
// ConvertGroupMod converts group mod.
//
func (c *DBCtl) ConvertGroupMod(mod *fibcapi.GroupMod) error {
	reID := mod.ReId

	switch e := mod.Entry.(type) {
	case *fibcapi.GroupMod_L2Iface:
		return c.ConvertL2IfaceGroup(reID, e.L2Iface)

	case *fibcapi.GroupMod_L3Unicast:
		return c.ConvertL3UnicastGroup(reID, e.L3Unicast)

	case *fibcapi.GroupMod_MplsIface:
		return c.ConvertMplsIfaceGroup(reID, e.MplsIface)

	case *fibcapi.GroupMod_MplsLabel:
		return nil

	default:
		return fmt.Errorf("Invalid flow mod. %s %v", reID, e)
	}
}

//
// ConvertL2Addrs converts l2addr list.
//
func (c *DBCtl) ConvertL2Addrs(dpID uint64, l2addrs []*fibcapi.L2Addr) (string, []*fibcapi.L2Addr, error) {
	reID, err := c.ConvertIDDPtoVM(dpID)
	if err != nil {
		return "", nil, err
	}

	vmAddrs := []*fibcapi.L2Addr{}
	for _, l2addr := range l2addrs {
		ok := c.PortMap().SelectByDP(dpID, l2addr.PortId, func(e *fibcdbm.PortEntry) {
			l2addr.PortId = e.VMPort.PortID
			l2addr.Ifname = e.Key.Ifname
			vmAddrs = append(vmAddrs, l2addr)
		})

		if !ok {
			c.log.Warnf("ConvertL2Addts: %s", err)
		}
	}

	return reID, vmAddrs, nil
}
