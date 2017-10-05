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

package ribctl

import (
	"fabricflow/fibc/api"
	"fabricflow/fibc/net"
	log "github.com/sirupsen/logrus"
	"time"
)

const FIBC_CONNECT_INTERVAL_MSEC = 1000 * time.Millisecond

type FIBController struct {
	active bool
	fibcon *fibcapi.FIBCon
	connCh chan bool
	recvCh chan fibcnet.Message
}

func NewFIBController(addr string) *FIBController {
	return &FIBController{
		active: false,
		fibcon: fibcapi.NewFIBCon(addr),
		connCh: make(chan bool),
		recvCh: make(chan fibcnet.Message),
	}
}

func (f *FIBController) Conn() <-chan bool {
	return f.connCh
}

func (f *FIBController) Recv() <-chan fibcnet.Message {
	return f.recvCh
}

func (f *FIBController) Start() {
	f.active = true

	go func() {
		for {
			if err := f.fibcon.Connect(); err == nil {
				f.Monitor()
			}
			if f.active != true {
				break
			}
			log.Debug("FIBController: Connectiong...")
			time.Sleep(FIBC_CONNECT_INTERVAL_MSEC)
		}
	}()
}

func (f *FIBController) Stop() {
	f.active = false
	f.fibcon.Close()
	time.Sleep(FIBC_CONNECT_INTERVAL_MSEC)
}

func (f *FIBController) Monitor() {
	f.connCh <- true // connected
	defer func() {
		f.connCh <- false // disconnected
	}()

	log.Infof("FIBController: Monitor START.")

	for {
		hdr, data, err := f.fibcon.Read()
		if err != nil {
			log.Errorf("FIBController: Monitor EXIT. Read error. %s", err)
			return
		}

		switch fibcapi.FFM(hdr.Type) {
		case fibcapi.FFM_PORT_STATUS:
			msg, err := fibcapi.NewPortStatusFromBytes(data)
			if err != nil {
				log.Errorf("FIBController: PortStatus error. %s", err)
				continue
			}
			f.recvCh <- msg

		case fibcapi.FFM_DP_STATUS:
			msg, err := fibcapi.NewDpStatusFromBytes(data)
			if err != nil {
				log.Errorf("FIBController: DpStatus error. %s", err)
				continue
			}
			f.recvCh <- msg

		default:
			log.Warnf("FIBController: Drop message. %v", hdr)
		}
	}
}

func (f *FIBController) Send(msg fibcnet.Message, xid uint32) error {
	return f.fibcon.Write(msg, xid)
}

type FIBCHandler interface {
	OnPortStatus(*fibcapi.PortStatus)
	OnDpStatus(*fibcapi.DpStatus)
}

func FIBCDispatch(msg fibcnet.Message, h FIBCHandler) {
	msgType := fibcapi.FFM(msg.Type())
	switch msgType {
	case fibcapi.FFM_PORT_STATUS:
		h.OnPortStatus(msg.(*fibcapi.PortStatus))
	case fibcapi.FFM_DP_STATUS:
		h.OnDpStatus(msg.(*fibcapi.DpStatus))
	default:
		log.Warnf("RIBController: Drop FIBC message. %v", msg)
	}
}
