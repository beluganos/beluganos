package netlink

import (
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

func LinkSetMulticastOn(link Link) error {
	return pkgHandle.LinkSetMulticastOn(link)
}

func (h *Handle) LinkSetMulticastOn(link Link) error {
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(unix.RTM_NEWLINK, unix.NLM_F_ACK)

	msg := nl.NewIfInfomsg(unix.AF_UNSPEC)
	msg.Change = unix.IFF_MULTICAST
	msg.Flags = unix.IFF_MULTICAST

	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(unix.NETLINK_ROUTE, 0)
	return err
}

func LinkSetMulticastOff(link Link) error {
	return pkgHandle.LinkSetMulticastOff(link)
}

func (h *Handle) LinkSetMulticastOff(link Link) error {
	base := link.Attrs()
	h.ensureIndex(base)
	req := h.newNetlinkRequest(unix.RTM_NEWLINK, unix.NLM_F_ACK)

	msg := nl.NewIfInfomsg(unix.AF_UNSPEC)
	msg.Change = unix.IFF_MULTICAST
	msg.Index = int32(base.Index)
	req.AddData(msg)

	_, err := req.Execute(unix.NETLINK_ROUTE, 0)
	return err
}
