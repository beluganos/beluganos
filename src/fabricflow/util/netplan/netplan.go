// -*- coding: utf-8 -*-

package netplan

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	NETWROK           = "network"
	VERSION           = "version"
	RENDERER          = "renderer"
	RENDERER_NETWORKD = "networkd"
	RENDERER_NETWORKM = "NetworkManager"
	ADDRESSES         = "addresses"
	DHCP4             = "dhcp4"
	DHCP6             = "dhcp6"
	ETHERNETS         = "ethernets"
	VLANS             = "vlans"
	VLANS_ID          = "id"
	VLANS_LINK        = "link"
	BONDS             = "bonds"
	BONDS_IFACES      = "interfaces"
	BONDS_PARAMS      = "parameters"
)

func NewConfig() map[interface{}]interface{} {
	return map[interface{}]interface{}{}
}

func WriteConfig(w io.Writer, c map[interface{}]interface{}) error {
	e := yaml.NewEncoder(w)
	defer e.Close()
	return e.Encode(c)
}

func ReadConfig(r io.Reader) (map[interface{}]interface{}, error) {
	d := yaml.NewDecoder(r)
	c := NewConfig()
	if err := d.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}

func NewEthernetPath(ifname string) []string {
	return []string{
		NETWROK,
		ETHERNETS,
		ifname,
	}
}

func NewVlanPath(ifname string, vid uint16) []string {
	if vid != 0 {
		ifname = fmt.Sprintf("%s.%d", ifname, vid)
	}

	return []string{
		NETWROK,
		VLANS,
		ifname,
	}
}

func NewBondPath(ifname string) []string {
	return []string{
		NETWROK,
		BONDS,
		ifname,
	}
}
