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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//
// DpsMsg dpid http message.
//
type DpsMsg struct {
	DpIds []uint64 `json:"dpids"`
}

func NewDpsMsg() *DpsMsg {
	return &DpsMsg{
		DpIds: []uint64{},
	}
}

//
// FIBHttpController is FIBC Client (HTTP protocol/ JSON format)
//
type FIBHttpController struct {
	baseUrl string

	log *log.Entry
}

//
// NewFIBHttpController returns new FIBHttpController.
//
func NewFIBHttpController(baseUrl string) *FIBHttpController {
	return &FIBHttpController{
		baseUrl: baseUrl,

		log: log.WithFields(log.Fields{"module": "FIBHttpController"}),
	}
}

//
// DpIds returns dpid list.
//
func (c *FIBHttpController) Dps() ([]uint64, error) {
	url := fmt.Sprintf("%s/fib/dps", c.baseUrl)

	c.log.Debugf("Dps: %s", url)

	res, err := HTTPGet(url)
	if err != nil {
		c.log.Errorf("Dps: HTTPGet error. %s", err)
		return nil, err
	}

	msg := NewDpsMsg()
	if err := json.NewDecoder(res).Decode(msg); err != nil {
		c.log.Errorf("Dps: decode error. %s", err)
		return nil, err
	}

	return msg.DpIds, nil
}

//
// PortStats returns port stats list for each ports.
//
func (c *FIBHttpController) PortStats(dpId uint64, names []string) (PortStats, error) {
	url := fmt.Sprintf("%s/fib/stats/port/%d", c.baseUrl, dpId)

	c.log.Debugf("PortStats: %s", url)

	res, err := HTTPGet(url)
	if err != nil {
		c.log.Errorf("PortStats: HttpGet error. %s %s", url, err)
		return nil, err
	}

	msgs := map[string]PortStats{}
	if err := json.NewDecoder(res).Decode(&msgs); err != nil {
		c.log.Errorf("PortStats: json.Decode error. %s", err)
		return nil, err
	}

	msg, ok := msgs[fmt.Sprintf("%d", dpId)]
	if !ok {
		c.log.Errorf("PortStats: dpid not found. %d", dpId)
		return nil, fmt.Errorf("dpid not found. %d", dpId)
	}

	msg.normalize()

	return msg, nil
}

//
// HTTPGet execute http get.
//
func HTTPGet(url string) (io.Reader, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	return res.Body, nil
}
