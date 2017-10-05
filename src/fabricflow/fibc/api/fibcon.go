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

package fibcapi

import (
	"fabricflow/fibc/net"
	"fmt"
	"net"
)

type FIBCon struct {
	addr string
	conn net.Conn
}

func NewFIBCon(addr string) *FIBCon {
	return &FIBCon{
		addr: addr,
		conn: nil,
	}
}

func (c *FIBCon) Connect() error {

	c.Close()

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *FIBCon) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *FIBCon) Read() (*fibcnet.Header, []byte, error) {
	if c.conn == nil {
		return nil, nil, fmt.Errorf("FIBCon: Read error. DISCONNECTED.")
	}
	return fibcnet.Read(c.conn)
}

func (c *FIBCon) Write(msg fibcnet.Message, xid uint32) error {
	if c.conn == nil {
		return fmt.Errorf("FIBCon: Write error. DISCONNECTED.")
	}
	return fibcnet.WriteMessage(c.conn, msg, xid)
}
