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
	"net"

	"github.com/vishvananda/netlink"
)

const (
	LINK_TYPE_DEVICE  = "device"
	LINK_TYPE_BRIDGE  = "bridge"
	LINK_TYPE_VLAN    = "vlan"
	LINK_TYPE_VXLAN   = "vxlan"
	LINK_TYPE_VTI     = "vti"
	LINK_TYPE_VTI6    = "vti6"
	LINK_TYPE_VETH    = "veth"
	LINK_TYPE_BOND    = "bond"
	LINK_TYPE_GENERIC = "generic"
	LINK_TYPE_IP4TUN  = "ipip"
	LINK_TYPE_IP6TUN  = "ip6tnl"
)

// LinkOperState
func ParseLinkOperState(s string) LinkOperState {
	if v, ok := LinkOperState_value[s]; ok {
		return LinkOperState(v)
	}

	return LinkOperState_OperUnknown
}

// LinkAttrs
func NewLinkAttrs() *LinkAttrs {
	return &LinkAttrs{
		HardwareAddr: net.HardwareAddr{},
	}
}

func (a *LinkAttrs) NetHardwareAddr() net.HardwareAddr {
	return net.HardwareAddr(a.HardwareAddr)
}

func LinkAttrsToNative(a *LinkAttrs, n *netlink.LinkAttrs) {
	n.Index = int(a.Index)
	n.MTU = int(a.Mtu)
	n.TxQLen = int(a.TxQLen)
	n.Name = a.Name
	n.HardwareAddr = a.HardwareAddr
	n.Flags = net.Flags(a.Flags)
	n.RawFlags = a.RawFlags
	n.ParentIndex = int(a.ParentIndex)
	n.MasterIndex = int(a.MasterIndex)
	n.Alias = a.Alias
	n.Promisc = int(a.Promisc)
	n.EncapType = a.EncapType
	n.OperState = netlink.LinkOperState(a.OperState)
	n.Slave = SlaveInfoToNative(a.GetSlaveInfo())
}

func LinkAttrsFromNative(a *LinkAttrs, n *netlink.LinkAttrs) {
	a.Index = int32(n.Index)
	a.Mtu = int32(n.MTU)
	a.TxQLen = int32(n.TxQLen)
	a.Name = n.Name
	a.HardwareAddr = n.HardwareAddr
	a.Flags = uint32(n.Flags)
	a.RawFlags = n.RawFlags
	a.ParentIndex = int32(n.ParentIndex)
	a.MasterIndex = int32(n.MasterIndex)
	a.Alias = n.Alias
	a.Promisc = int32(n.Promisc)
	a.EncapType = n.EncapType
	a.OperState = LinkOperState(n.OperState)
	a.SlaveInfo = NewSlaveInfoFromNative(n.Slave)
}

// Link (Generic)
func NewGenericLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_GENERIC,
		LinkAttrs: NewGenericLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}

}
func NewGenericLinkAttrs() *Link_Generic {
	return &Link_Generic{
		Generic: &GenericLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
		},
	}
}

func NewGenericLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	a := NewGenericLinkAttrs()
	LinkAttrsFromNative(a.Generic.GetLinkAttrs(), ln.Attrs())

	return a
}

func GenericLinkToNative(ln *Link) netlink.Link {
	n := &netlink.GenericLink{}
	a := ln.GetGeneric()
	LinkAttrsToNative(a.GetLinkAttrs(), n.Attrs())

	n.LinkType = LINK_TYPE_GENERIC

	return n
}

// Link (Device)
func NewDeviceLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_DEVICE,
		LinkAttrs: NewDeviceLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}
func NewDeviceLinkAttrs() *Link_Device {
	return &Link_Device{
		Device: &DeviceLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
		},
	}
}

func NewDeviceLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	a := NewDeviceLinkAttrs()
	LinkAttrsFromNative(a.Device.GetLinkAttrs(), ln.Attrs())

	return a
}

func DeviceLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Device{}
	a := ln.GetDevice()
	LinkAttrsToNative(a.GetLinkAttrs(), n.Attrs())

	return n
}

// Link (Bridge)
func NewBridgeLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_BRIDGE,
		LinkAttrs: NewBridgeLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewBridgeLinkAttrs() *Link_Bridge {
	return &Link_Bridge{
		Bridge: &BridgeLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
		},
	}
}

func NewBridgeLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	a := NewBridgeLinkAttrs()
	LinkAttrsFromNative(a.Bridge.GetLinkAttrs(), ln.Attrs())

	n := ln.(*netlink.Bridge)
	a.Bridge.MulticastSnooping = func() bool {
		if mcSnoop := n.MulticastSnooping; mcSnoop != nil {
			return *mcSnoop
		}
		return false
	}()

	a.Bridge.HelloTime = func() uint32 {
		if helloTime := n.HelloTime; helloTime != nil {
			return *helloTime
		}
		return 0
	}()

	a.Bridge.VlanFiltering = func() bool {
		if vlanFiltering := n.VlanFiltering; vlanFiltering != nil {
			return *vlanFiltering
		}
		return false
	}()

	return a
}

func BridgeLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Bridge{}
	a := ln.GetBridge()
	LinkAttrsToNative(a.GetLinkAttrs(), n.Attrs())

	mcSnoop := a.MulticastSnooping
	n.MulticastSnooping = &mcSnoop

	helloTime := a.HelloTime
	n.HelloTime = &helloTime

	vlanFiltering := a.VlanFiltering
	n.VlanFiltering = &vlanFiltering

	return n
}

// Link (VLAN)
func NewVlanLink(nid uint8, lnId uint16, vlanId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_VLAN,
		LinkAttrs: NewVlanLinkAttrs(vlanId),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewVlanLinkAttrs(vlanId uint16) *Link_Vlan {
	return &Link_Vlan{
		Vlan: &VlanLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
			VlanId:    int32(vlanId),
		},
	}
}

func NewVlanLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	n := ln.(*netlink.Vlan)
	a := NewVlanLinkAttrs(uint16(n.VlanId))
	LinkAttrsFromNative(a.Vlan.GetLinkAttrs(), ln.Attrs())

	return a
}

func VlanLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Vlan{}
	a := ln.GetVlan()
	LinkAttrsToNative(a.GetLinkAttrs(), n.Attrs())

	n.VlanId = int(a.VlanId)

	return n
}

// Link (VxLAN)
func NewVxlanLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_VXLAN,
		LinkAttrs: NewVxlanLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewVxlanLinkAttrs() *Link_Vxlan {
	return &Link_Vxlan{
		Vxlan: &VxlanLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
			SrcAddr:   net.IP{},
			Group:     net.IP{},
		},
	}
}

func NewVxlanLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	attr := NewVxlanLinkAttrs()
	LinkAttrsFromNative(attr.Vxlan.GetLinkAttrs(), ln.Attrs())

	n := ln.(*netlink.Vxlan)
	a := attr.Vxlan
	a.VxlanId = int32(n.VxlanId)
	a.VtepDevIndex = int32(n.VtepDevIndex)
	a.SrcAddr = n.SrcAddr
	a.Group = n.Group
	a.Ttl = int32(n.TTL)
	a.Tos = int32(n.TOS)
	a.Learning = n.Learning
	a.Proxy = n.Proxy
	a.Rsc = n.RSC
	a.L2Miss = n.L2miss
	a.L3Miss = n.L3miss
	a.UdpCSum = n.UDPCSum
	a.NoAge = n.NoAge
	a.Gbp = n.GBP
	a.Age = int32(n.Age)
	a.Limit = int32(n.Limit)
	a.Port = int32(n.Port)
	a.PortLow = int32(n.PortLow)
	a.PortHigh = int32(n.PortHigh)

	return attr
}

func VxlanLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Vxlan{}
	d := ln.GetVxlan()
	LinkAttrsToNative(d.GetLinkAttrs(), n.Attrs())

	n.VxlanId = int(d.VxlanId)
	n.VtepDevIndex = int(d.VtepDevIndex)
	n.SrcAddr = net.IP(d.SrcAddr)
	n.Group = net.IP(d.Group)
	n.TTL = int(d.Ttl)
	n.TOS = int(d.Tos)
	n.Learning = d.Learning
	n.Proxy = d.Proxy
	n.RSC = d.Rsc
	n.L2miss = d.L2Miss
	n.L3miss = d.L3Miss
	n.UDPCSum = d.UdpCSum
	n.NoAge = d.NoAge
	n.GBP = d.Gbp
	n.Age = int(d.Age)
	n.Limit = int(d.Limit)
	n.Port = int(d.Port)
	n.PortLow = int(d.PortLow)
	n.PortHigh = int(d.PortHigh)

	return n
}

// Link (Vti)
func NewVtiLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_VTI,
		LinkAttrs: NewVtiLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewVtiLinkAttrs() *Link_Vti {
	return &Link_Vti{
		Vti: &VtiLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
			Local:     net.IP{},
			Remote:    net.IP{},
		},
	}
}

func NewVtiLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	attr := NewVtiLinkAttrs()
	LinkAttrsFromNative(attr.Vti.GetLinkAttrs(), ln.Attrs())

	n := ln.(*netlink.Vti)
	a := attr.Vti
	a.IKey = n.IKey
	a.OKey = n.OKey
	a.Link = n.Link
	a.Local = n.Local
	a.Remote = n.Remote

	return attr
}

func VtiLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Vti{}
	d := ln.GetVti()
	LinkAttrsToNative(d.GetLinkAttrs(), n.Attrs())

	n.IKey = d.IKey
	n.OKey = d.OKey
	n.Link = d.Link
	n.Local = net.IP(d.Local)
	n.Remote = net.IP(d.Remote)

	return n
}

// Link (Veth)
func NewVethLink(nid uint8, lnId uint16, peerName string) *Link {
	return &Link{
		Type:      LINK_TYPE_VETH,
		LinkAttrs: NewVethLinkAttrs(peerName),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewVethLinkAttrs(peerName string) *Link_Veth {
	return &Link_Veth{
		Veth: &VethLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
			PeerName:  peerName,
		},
	}
}

func NewVethLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	n := ln.(*netlink.Veth)
	a := NewVethLinkAttrs(n.PeerName)
	LinkAttrsFromNative(a.Veth.GetLinkAttrs(), ln.Attrs())

	return a
}

func VethLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Veth{}
	d := ln.GetVeth()
	LinkAttrsToNative(d.GetLinkAttrs(), n.Attrs())

	n.PeerName = d.PeerName

	return n
}

// Link (BondAdInfo)
func BondAdInfoToNative(a *BondAdInfo) *netlink.BondAdInfo {
	if a == nil {
		return nil
	}

	n := &netlink.BondAdInfo{}

	n.AggregatorId = int(a.AggregatorId)
	n.NumPorts = int(a.NumPorts)
	n.ActorKey = int(a.ActorKey)
	n.PartnerKey = int(a.PartnerKey)
	n.PartnerMac = net.HardwareAddr(a.PartnerMac)

	return n
}

func BondAdInfoFromNative(n *netlink.BondAdInfo) *BondAdInfo {
	if n == nil {
		return nil
	}

	a := &BondAdInfo{}

	a.AggregatorId = int32(n.AggregatorId)
	a.NumPorts = int32(n.NumPorts)
	a.ActorKey = int32(n.ActorKey)
	a.PartnerKey = int32(n.PartnerKey)
	a.PartnerMac = n.PartnerMac

	return a
}

func NewBondAdInfo() *BondAdInfo {
	return &BondAdInfo{
		PartnerMac: net.HardwareAddr{},
	}
}

// Link (Bond)
func NewBondLink(nid uint8, lnId uint16) *Link {
	return &Link{
		Type:      LINK_TYPE_BOND,
		LinkAttrs: NewBondLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewBondLinkAttrs() *Link_Bond {
	return &Link_Bond{
		Bond: &BondLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
			AdInfo:    NewBondAdInfo(),
		},
	}
}

func NewBondLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	attr := NewBondLinkAttrs()
	LinkAttrsFromNative(attr.Bond.GetLinkAttrs(), ln.Attrs())

	n := ln.(*netlink.Bond)
	a := attr.Bond
	a.Mode = BondMode(n.Mode)
	a.ActiveSlave = int32(n.ActiveSlave)
	a.Miimon = int32(n.Miimon)
	a.UpDelay = int32(n.UpDelay)
	a.DownDelay = int32(n.DownDelay)
	a.UseCarrier = int32(n.UseCarrier)
	a.ArpInterval = int32(n.ArpInterval)
	a.ArpIpTargets = IPsToBytes(n.ArpIpTargets)
	a.ArpValidate = BondArpValidate(n.ArpValidate)
	a.ArpAllTargets = BondArpAllTargets(n.ArpAllTargets)
	a.Primary = int32(n.Primary)
	a.PrimaryReselect = BondPrimaryReselect(n.PrimaryReselect)
	a.FailOverMac = BondFailOverMac(n.FailOverMac)
	a.XmitHashPolicy = BondXmitHashPolicy(n.XmitHashPolicy)
	a.ResendIgmp = int32(n.ResendIgmp)
	a.NumPeerNotif = int32(n.NumPeerNotif)
	a.AllSlavesActive = int32(n.AllSlavesActive)
	a.MinLinks = int32(n.MinLinks)
	a.LpInterval = int32(n.LpInterval)
	a.PackersPerSlave = int32(n.PackersPerSlave)
	a.LacpRate = BondLacpRate(n.LacpRate)
	a.AdSelect = BondAdSelect(n.AdSelect)
	a.AdInfo = BondAdInfoFromNative(n.AdInfo)

	return attr
}

func BondLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Bond{}
	a := ln.GetBond()
	LinkAttrsToNative(a.GetLinkAttrs(), n.Attrs())

	n.Mode = netlink.BondMode(a.Mode)
	n.ActiveSlave = int(a.ActiveSlave)
	n.Miimon = int(a.Miimon)
	n.UpDelay = int(a.UpDelay)
	n.DownDelay = int(a.DownDelay)
	n.UseCarrier = int(a.UseCarrier)
	n.ArpInterval = int(a.ArpInterval)
	n.ArpIpTargets = BytesToIPs(a.ArpIpTargets)
	n.ArpValidate = netlink.BondArpValidate(a.ArpValidate)
	n.ArpAllTargets = netlink.BondArpAllTargets(a.ArpAllTargets)
	n.Primary = int(a.Primary)
	n.PrimaryReselect = netlink.BondPrimaryReselect(a.PrimaryReselect)
	n.FailOverMac = netlink.BondFailOverMac(a.FailOverMac)
	n.XmitHashPolicy = netlink.BondXmitHashPolicy(a.XmitHashPolicy)
	n.ResendIgmp = int(a.ResendIgmp)
	n.NumPeerNotif = int(a.NumPeerNotif)
	n.AllSlavesActive = int(a.AllSlavesActive)
	n.MinLinks = int(a.MinLinks)
	n.LpInterval = int(a.LpInterval)
	n.PackersPerSlave = int(a.PackersPerSlave)
	n.LacpRate = netlink.BondLacpRate(a.LacpRate)
	n.AdSelect = netlink.BondAdSelect(a.AdSelect)
	n.AdInfo = BondAdInfoToNative(a.AdInfo)

	return n
}

// Link (Iptun)
func NewIptunLink(nid uint8, lnId uint8, typ string) *Link {
	return &Link{
		Type:      typ,
		LinkAttrs: NewIptunLinkAttrs(),
		NId:       uint32(nid),
		LnId:      uint32(lnId),
	}
}

func NewIptunLinkAttrs() *Link_Iptun {
	return &Link_Iptun{
		Iptun: &IptunLinkAttrs{
			LinkAttrs: NewLinkAttrs(),
		},
	}
}

func NewIptunLinkAttrsFromNative(ln netlink.Link) isLink_LinkAttrs {
	attrs := NewIptunLinkAttrs()
	a := attrs.Iptun
	LinkAttrsFromNative(a.GetLinkAttrs(), ln.Attrs())

	n := ln.(*netlink.Iptun)
	a.Ttl = uint32(n.Ttl)
	a.Tos = uint32(n.Tos)
	a.PMtuDisc = uint32(n.PMtuDisc)
	a.Link = uint32(n.Link)
	a.Local = n.Local
	a.Remote = n.Remote
	a.EncapSport = uint32(n.EncapSport)
	a.EncapDport = uint32(n.EncapDport)
	a.EncapType = uint32(n.EncapType)
	a.EncapFlags = uint32(n.EncapFlags)
	a.FlowBased = n.FlowBased

	return attrs
}

func IptunLinkToNative(ln *Link) netlink.Link {
	n := &netlink.Iptun{}
	d := ln.GetIptun()
	LinkAttrsToNative(d.GetLinkAttrs(), n.Attrs())

	n.Ttl = uint8(d.Ttl)
	n.Tos = uint8(d.Tos)
	n.PMtuDisc = uint8(d.PMtuDisc)
	n.Link = uint32(d.Link)
	n.Local = net.IP(d.Local)
	n.Remote = net.IP(d.Remote)
	n.EncapSport = uint16(d.EncapSport)
	n.EncapDport = uint16(d.EncapDport)
	n.EncapType = uint16(d.EncapType)
	n.EncapFlags = uint16(d.EncapFlags)
	n.FlowBased = d.FlowBased

	return n
}

// Link
var linkToNativeFuncs = map[string]func(*Link) netlink.Link{
	LINK_TYPE_DEVICE:  DeviceLinkToNative,
	LINK_TYPE_BRIDGE:  BridgeLinkToNative,
	LINK_TYPE_VLAN:    VlanLinkToNative,
	LINK_TYPE_VXLAN:   VxlanLinkToNative,
	LINK_TYPE_VTI:     VtiLinkToNative,
	LINK_TYPE_VETH:    VethLinkToNative,
	LINK_TYPE_BOND:    BondLinkToNative,
	LINK_TYPE_GENERIC: GenericLinkToNative,
	LINK_TYPE_IP4TUN:  IptunLinkToNative,
	LINK_TYPE_IP6TUN:  IptunLinkToNative,
}

var linkFromNativeFuncs = map[string]func(netlink.Link) isLink_LinkAttrs{
	LINK_TYPE_DEVICE:  NewDeviceLinkAttrsFromNative,
	LINK_TYPE_BRIDGE:  NewBridgeLinkAttrsFromNative,
	LINK_TYPE_VLAN:    NewVlanLinkAttrsFromNative,
	LINK_TYPE_VXLAN:   NewVxlanLinkAttrsFromNative,
	LINK_TYPE_VTI:     NewVtiLinkAttrsFromNative,
	LINK_TYPE_VETH:    NewVethLinkAttrsFromNative,
	LINK_TYPE_BOND:    NewBondLinkAttrsFromNative,
	LINK_TYPE_GENERIC: NewGenericLinkAttrsFromNative,
	LINK_TYPE_IP4TUN:  NewIptunLinkAttrsFromNative,
	LINK_TYPE_IP6TUN:  NewIptunLinkAttrsFromNative,
}

func (ln *Link) ToNetlink() netlink.Link {
	f, ok := linkToNativeFuncs[ln.Type]
	if !ok {
		f = GenericLinkToNative
	}
	return f(ln)
}

func (ln *Link) ToNative() *nlamsg.Link {
	return &nlamsg.Link{
		Link: ln.ToNetlink(),
		LnId: uint16(ln.LnId),
		NId:  uint8(ln.NId),
	}
}

func NewLinkFromNative(ln *nlamsg.Link) *Link {
	f, ok := linkFromNativeFuncs[ln.Link.Type()]
	if !ok {
		f = NewGenericLinkAttrsFromNative
	}

	return &Link{
		Type:      ln.Type(),
		LinkAttrs: f(ln.Link),
		NId:       uint32(ln.NId),
		LnId:      uint32(ln.LnId),
	}
}

//
// Link (Key)
//
func (k *LinkKey) ToNative() *nladbm.LinkKey {
	return nladbm.NewLinkKey(uint8(k.NId), int(k.Index))
}

func NewLinkKeyFromNative(n *nladbm.LinkKey) *LinkKey {
	return &LinkKey{
		NId:   uint32(n.NId),
		Index: int32(n.Index),
	}
}

//
// Links
//
func NewGetLinksRequest(nid uint8) *GetLinksRequest {
	return &GetLinksRequest{
		NId: uint32(nid),
	}
}
