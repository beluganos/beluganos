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

package main

import (
	"net"
	"time"

	"github.com/mdlayher/ndp"
	log "github.com/sirupsen/logrus"
)

type Worker interface {
	Start() error
	Stop()
}

type worker struct {
	ifname     string
	conn       *ndp.Conn
	nsInterval time.Duration
	naTable    map[string]*ndp.NeighborAdvertisement
	log        *log.Entry
}

func NewWorker(ifname string, nsInterval time.Duration) Worker {
	return &worker{
		ifname:     ifname,
		nsInterval: nsInterval,
		naTable:    map[string]*ndp.NeighborAdvertisement{},
		log: log.WithFields(log.Fields{
			"ifname": ifname,
		}),
	}
}

func (w *worker) cleanTable() {
	w.naTable = map[string]*ndp.NeighborAdvertisement{}
}

func (w *worker) Start() error {
	ifi, err := net.InterfaceByName(w.ifname)
	if err != nil {
		return err
	}

	conn, ip, err := ndp.Dial(ifi, ndp.LinkLocal)
	if err != nil {
		return err
	}

	w.log.Debugf("Dial %v -> %s", ifi, ip)

	w.conn = conn

	ch := make(chan *ndp.NeighborAdvertisement)
	go w.Recv(ch)
	go w.Serve(ch)

	return nil
}

func (w *worker) Serve(ch <-chan *ndp.NeighborAdvertisement) {

	w.log.Debugf("Serve start. %s", w.ifname)

	ticker := time.NewTicker(w.nsInterval)
	defer ticker.Stop()

	for {
		select {
		case m := <-ch:
			if m == nil {
				w.log.Debugf("Serve exit. %s", w.ifname)
				return
			}

			if m.TargetAddress.IsGlobalUnicast() {
				w.naTable[m.TargetAddress.String()] = m
				w.log.Debugf("Serve add. %s NA:%v", w.ifname, m)
			}

		case <-ticker.C:
			for _, m := range w.naTable {
				w.log.Debugf("Serve send. %s NA:%v", w.ifname, m)

				ns := &ndp.NeighborSolicitation{
					TargetAddress: m.TargetAddress,
				}

				if err := w.conn.WriteTo(ns, nil, m.TargetAddress); err != nil {
					w.log.Errorf("Serve write error. %s", err)
				}

				w.log.Debugf("Serve send. %s NS:%v", w.ifname, m)
			}

			w.cleanTable()
		}
	}
}

func (w *worker) Recv(ch chan<- *ndp.NeighborAdvertisement) {
	defer close(ch)

	for {
		msg, cmsg, ip, err := w.conn.ReadFrom()
		if err != nil {
			w.log.Warnf("Recv %s %s", err, w.ifname)
			break
		}

		w.log.Debugf("Recv %s from %s", w.ifname, ip)
		w.log.Debugf("Recv %s %s", msg, cmsg)

		switch m := msg.(type) {
		case *ndp.NeighborAdvertisement:
			w.log.Debugf("Recv %s %s %v", w.ifname, m.Type(), m)
			ch <- m

		default:
			// w.log.Debugf("Recv %s %s %s", w.ifname, m.Type(), m)
		}
	}

	w.log.Debugf("Recv exit. %s", w.ifname)
}

func (w *worker) Stop() {
	w.conn.Close()
}
