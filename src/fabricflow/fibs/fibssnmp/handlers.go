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
	"fmt"

	lib "fabricflow/fibs/fibslib"

	log "github.com/sirupsen/logrus"
)

//
// PortStatsHandler is generic port stats handler.
//
type PortStatsHandler struct {
	BaseOid string
	Name    string
	Type    lib.SnmpType
	Datas   StatsDatas
}

//
// NewPortStatsHandler returns new instance.
//
func NewPortStatsHandler(baseOid, name string, t lib.SnmpType, datas StatsDatas) *PortStatsHandler {
	return &PortStatsHandler{
		BaseOid: baseOid,
		Name:    name,
		Type:    t,
		Datas:   datas,
	}
}

//
// NewPortStatsHandlerFromConfig returns new instance from config.
//
func NewPortStatsHandlerFromConfig(cfg *HandlerConfig, datas StatsDatas) *PortStatsHandler {
	return NewPortStatsHandler(cfg.Oid, cfg.Name, lib.SnmpType(cfg.Type), datas)
}

//
// String returns description
//
func (h *PortStatsHandler) String() string {
	return fmt.Sprintf("'%s', %s, %s, PortStats", h.BaseOid, h.Name, h.Type)
}

//
// newSnmpReply returns new snmp reply instance.
//
func (h *PortStatsHandler) newSnmpReply(ps PortStats) *SnmpReply {

	index, _ := ps.PortNo()
	value, ok := ps[h.Name]

	if !ok {
		value = ""
	}

	return NewSnmpReply(
		fmt.Sprintf("%s.%d", h.BaseOid, index),
		h.Type,
		value,
	)
}

//
// Oid returns oid.
//
func (h *PortStatsHandler) Oid() string {
	return h.BaseOid
}

//
// Get proces get request.
//
func (h *PortStatsHandler) Get(oid string) *SnmpReply {
	subOid := lib.ParseOID(oid[len(h.BaseOid):])
	log.Debugf("PortStatsHandler.Get: name = '%s', suboid = '%v'", h.Name, subOid)

	index, ok := GetIndexOfOID(subOid)
	if !ok {
		return nil
	}

	ps, ok := h.Datas.PortStatsList().Get(index)
	if !ok {
		return nil
	}

	return h.newSnmpReply(ps)
}

//
// GetNext process getnext request.
//
func (h *PortStatsHandler) GetNext(oid string) *SnmpReply {
	subOid := lib.ParseOID(oid[len(h.BaseOid):])
	log.Debugf("IFNamesHandler.GetNext: name = '%s', suboid = '%v'", h.Name, subOid)

	nextIndex, ok := GetNextIndexOfOID(subOid, 0)
	if !ok {
		return nil
	}

	ps, ok := h.Datas.PortStatsList().GetNext(nextIndex)
	if !ok {
		return nil
	}

	return h.newSnmpReply(ps)
}
