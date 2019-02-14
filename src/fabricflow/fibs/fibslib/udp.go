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

package fibslib

import (
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	UDP_BUFFER_SIZE = 4096
)

//
// UDPServerCallback is callback function
//
type UDPServerCallback func(buf []byte, client *net.UDPAddr, conn *net.UDPConn)

func StartUDPServer(laddr *net.UDPAddr, cb UDPServerCallback) error {
	log.Debugf("StartUDPServer %s", laddr)

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Errorf("UDPServer Listen error. %s", err)
		return err
	}

	log.Infof("StartUDPServer START")

	buf := make([]byte, UDP_BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("UDPServer Read error. %s", err)
			break
		}

		log.Debugf("UDPServer Read size:%d, addr:%s", n, addr)

		cb(buf[:n], addr, conn)
	}

	return nil
}
