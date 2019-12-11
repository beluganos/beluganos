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
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

const (
	OptionCodeRouter          = 3
	OptionCodeDefaultURL      = 114
	OptionCodeBeluganosKVMURL = 250
	OptionCodeBeluganosZTPURL = 251
)

const (
	OptionNameRouter          = "router"
	OptionNameDefaultURL      = "default-url"
	OptionNameBeluganosKVMURL = "beluganos-kvm-url"
	OptionNameBeluganosZTPURL = "beluganos-ztp-url"
)

type OptionDecoder func([]byte) (string, error)

type IPv4Option struct {
	Code   uint8
	Decode OptionDecoder
}

func (e IPv4Option) OptionCode() dhcpv4.OptionCode {
	return dhcpv4.GenericOptionCode(e.Code)
}

func decodeIPv4Option(b []byte) (string, error) {
	ip := net.IP(b)
	ip = ip.To4()
	if ip == nil {
		return "", fmt.Errorf("Invalid ip addr.")
	}

	return ip.String(), nil
}

func decodeIPv4ListOption(b []byte) (string, error) {
	var pos int
	var total = len(b)
	addrs := []string{}

FOR_LOOP:
	for {
		if total-pos < 4 {
			break FOR_LOOP
		}
		addr, err := decodeIPv4Option(b[pos : pos+4])
		if err != nil {
			return "", err
		}

		pos += 4
		addrs = append(addrs, addr)
	}

	return strings.Join(addrs, ","), nil
}

func decodeStringOption(b []byte) (string, error) {
	return string(b), nil
}

var dhcpIPv4Options = map[string]*IPv4Option{
	OptionNameRouter:          &IPv4Option{Code: OptionCodeRouter, Decode: decodeIPv4ListOption},
	OptionNameDefaultURL:      &IPv4Option{Code: OptionCodeDefaultURL, Decode: decodeStringOption},
	OptionNameBeluganosKVMURL: &IPv4Option{Code: OptionCodeBeluganosKVMURL, Decode: decodeStringOption},
	OptionNameBeluganosZTPURL: &IPv4Option{Code: OptionCodeBeluganosZTPURL, Decode: decodeStringOption},
}

func ListIPv4Options(f func(string, *IPv4Option)) {
	for name, option := range dhcpIPv4Options {
		f(name, option)
	}
}

func GetIPv4Option(name string) (*IPv4Option, error) {
	if e, ok := dhcpIPv4Options[name]; ok {
		return e, nil
	}

	return nil, fmt.Errorf("option not found. %s", name)
}

func GetIPv4OptionByCode(code uint8) (*IPv4Option, error) {
	for _, e := range dhcpIPv4Options {
		if e.Code == code {
			return e, nil
		}
	}

	return nil, fmt.Errorf("option code not found. %d", code)
}
