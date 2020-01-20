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
	fibcapi "fabricflow/fibc/api"
	fibcnet "fabricflow/fibc/net"
	"fabricflow/fibc/pkgs/fibcdbm"
	"fmt"
	"log/syslog"
	"sync/atomic"
)

//
// DBMpWaiter is waiter for multipart message,
//
type DBMpWaiter struct {
	*fibcdbm.SimpleWaiter
	Reply *fibcapi.FFMultipart_Reply
}

//
// Set sets reply and close wait channel.
//
func (w *DBMpWaiter) Set(k interface{}, v interface{}) {
	if reply, ok := v.(*fibcapi.FFMultipart_Reply); ok {
		w.Reply = reply
		w.SimpleWaiter.Set(k, v)
		return
	}

	w.SimpleWaiter.SetError(fmt.Errorf("Invalid reply. %v", v))
}

//
// NewDBMpWaiter returns new DBMpWaiter
//
func NewDBMpWaiter() *DBMpWaiter {
	return &DBMpWaiter{
		SimpleWaiter: fibcdbm.NewSimpleWaiter(),
	}
}

//
// OAMWaiter is waiter for oam(audit route cnt) message
//
type OAMWaiter struct {
	*fibcdbm.SimpleWaiter
	Request *fibcapi.OAM_Request
	RestNum int32
	VMReply map[string]*fibcapi.OAM_Reply // key: re_id
	VSReply map[uint64]*fibcapi.OAM_Reply // key: vs_id
	DPReply map[uint64]*fibcapi.OAM_Reply // key: dp_id
}

func NewOAMWaiter(request *fibcapi.OAM_Request) *OAMWaiter {
	return &OAMWaiter{
		SimpleWaiter: fibcdbm.NewSimpleWaiter(),

		Request: request,
		VMReply: map[string]*fibcapi.OAM_Reply{},
		VSReply: map[uint64]*fibcapi.OAM_Reply{},
		DPReply: map[uint64]*fibcapi.OAM_Reply{},
	}
}

func (w *OAMWaiter) SetREID(reID string) {
	w.VMReply[reID] = nil
	atomic.AddInt32(&w.RestNum, 1)
}

func (w *OAMWaiter) SetVSID(vsID uint64) {
	w.VSReply[vsID] = nil
	atomic.AddInt32(&w.RestNum, 1)
}

func (w *OAMWaiter) SetDPID(dpID uint64) {
	w.DPReply[dpID] = nil
	atomic.AddInt32(&w.RestNum, 1)
}

//
// Set sets reply key/val.
//
// k: "dp" or "vs" or re_id
// v: *fibcapi.OAM_Reply(AuditRouteCntReply)
//
func (w *OAMWaiter) Set(k interface{}, v interface{}) {
	key, ok := k.(string)
	if !ok {
		w.SimpleWaiter.SetError(fmt.Errorf("Invalid reply key. %v %v", k, v))
		return
	}
	reply, ok := v.(*fibcapi.OAM_Reply)
	if !ok {
		w.SimpleWaiter.SetError(fmt.Errorf("Invalid reply val. %v %v", k, v))
		return
	}

	switch key {
	case "dp":
		w.DPReply[reply.DpId] = reply
	case "vs":
		w.VSReply[reply.DpId] = reply
	default:
		// key is re_id.
		w.VMReply[key] = reply
	}

	if restNum := atomic.AddInt32(&w.RestNum, -1); restNum <= 0 {
		w.SimpleWaiter.Set(k, v)
	}
}

func (w *OAMWaiter) Check() error {
	return fibcapi.DispatchOAMRequest(nil, w.Request, w)
}

func (w *OAMWaiter) FIBCOAMAuditRouteCntRequest(hdr *fibcnet.Header, oam *fibcapi.OAM_Request, audit *fibcapi.OAM_AuditRouteCntRequest) error {
	var (
		vmSum uint64
		dpSum uint64
	)

	for _, reply := range w.VMReply {
		vmSum += reply.GetAuditRouteCnt().GetCount()
	}

	for _, reply := range w.DPReply {
		dpSum += reply.GetAuditRouteCnt().GetCount()
	}

	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, "fibcd")
	if err != nil {
		return err
	}
	defer logger.Close()

	if vmSum != dpSum {
		logger.Err(fmt.Sprintf("audit route cnt: error. vm:%d dp:%d", vmSum, dpSum))
	} else {
		logger.Info(fmt.Sprintf("audit route cnt: ok. vm:%d dp:%d", vmSum, dpSum))
	}

	return nil
}
