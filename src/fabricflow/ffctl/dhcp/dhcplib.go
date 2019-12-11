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

package dhcplib

import (
	"fmt"
	"net"
	"regexp"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/interfaces"
	log "github.com/sirupsen/logrus"
)

func IPv4Exchanges(client *client4.Client, optCode dhcpv4.OptionCode, ifnames []string) (*dhcpv4.DHCPv4, error) {
	var (
		err error
		msg *dhcpv4.DHCPv4
	)

	for _, ifname := range ifnames {
		if msg, err = IPv4Exchange(client, optCode, ifname); err == nil {
			log.Debugf("dhcpIPv4Exchange ok. %s", ifname)
			return msg, nil
		}

		log.Debugf("dhcpIPv4Exchange error. %s %s", ifname, err)
	}

	return nil, err
}

func IPv4Exchange(client *client4.Client, optCode dhcpv4.OptionCode, ifname string) (*dhcpv4.DHCPv4, error) {
	opts := []dhcpv4.Modifier{
		dhcpv4.WithRequestedOptions(optCode),
	}

	msgs, err := client.Exchange(ifname, opts...)
	if err != nil {
		log.Debugf("DHCPv4.Exchange error. %s", err)
		return nil, err
	}

	if len(msgs) == 0 {
		log.Debugf("DHCPv4.Exchange bad messages.")
		return nil, fmt.Errorf("DHCPv4.Exchange bad messages.")
	}

	lastMsg := msgs[len(msgs)-1]
	if lastMsg.OpCode != dhcpv4.OpcodeBootReply {
		log.Debugf("DHCPv4.Exchange failed.")
		return nil, fmt.Errorf("DHCPv4.Exchange failed.")
	}

	return lastMsg, nil
}

func IfNameList(patterns []string, prilist []string, blacklist []string) []string {
	blackmap := map[string]struct{}{}
	for _, ifname := range blacklist {
		blackmap[ifname] = struct{}{}
	}

	primap := map[string]struct{}{}
	for _, ifname := range prilist {
		primap[ifname] = struct{}{}
	}

	relist := []*regexp.Regexp{}
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			relist = append(relist, re)
		} else {
			log.Debugf("linkList: regexp compile error. %s", pattern)
		}
	}

	ifaces, err := interfaces.GetInterfacesFunc(func(iface net.Interface) bool {
		if iface.Flags&net.FlagLoopback != 0 {
			log.Debugf("linkList: %s is loopback.", iface.Name)
			return false
		}
		if _, exist := blackmap[iface.Name]; exist {
			log.Debugf("linkList: %s is blacklist.", iface.Name)
			return false
		}
		for _, re := range relist {
			if ok := re.MatchString(iface.Name); ok {
				log.Debugf("linkList: %s match pattern. %s", iface.Name, re)
				return true
			}
			log.Debugf("linkList: %s unmatch pattern. %s", iface.Name, re)
		}

		return false
	})

	if err != nil {
		return []string{}
	}

	ifnames := []string{}
	for _, iface := range ifaces {
		if _, ok := primap[iface.Name]; ok {
			ifnames = append([]string{iface.Name}, ifnames...)
		} else {
			ifnames = append(ifnames, iface.Name)
		}
	}

	return ifnames
}
