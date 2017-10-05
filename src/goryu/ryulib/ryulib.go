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

package ryulib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goryu/encoding"
	"goryu/ofproto"
	"io"
	"net/http"
)

func HttpGet(url string) (io.Reader, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	return res.Body, nil
}

func HttpPost(url string, body []byte) error {
	r := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", r)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Error Response. %d", res.StatusCode)
	}

	return nil
}

type RyuClient struct {
	url string
}

func NewClient(url string) *RyuClient {
	return &RyuClient{
		url: url,
	}
}

func (c *RyuClient) GetSwitches() ([]int, error) {
	r, err := HttpGet(fmt.Sprintf("%s/%s", c.url, "switches"))
	if err != nil {
		return nil, err
	}

	switches := []int{}
	d := json.NewDecoder(r)
	return switches, d.Decode(&switches)
}

func (c *RyuClient) GetDesc(dpid int) (*ofproto.Desc, error) {
	r, err := HttpGet(fmt.Sprintf("%s/%s/%d", c.url, "desc", dpid))
	if err != nil {
		return nil, err
	}

	descs := make(map[int]ofproto.Desc)
	d := json.NewDecoder(r)
	if err := d.Decode(&descs); err != nil {
		return nil, err
	}

	for id, desc := range descs {
		if id == dpid {
			return &desc, nil
		}
	}

	return nil, fmt.Errorf("Internal error")
}

func (c *RyuClient) GetFlow(dpid int) ([]*ofproto.FlowEntry, error) {
	r, err := HttpGet(fmt.Sprintf("%s/%s/%d", c.url, "flow", dpid))
	if err != nil {
		return nil, err
	}

	flows := make(map[int][]map[string]interface{})
	d := json.NewDecoder(r)
	if err := d.Decode(&flows); err != nil {
		return nil, err
	}

	for id, flow := range flows {
		if id == dpid {
			return ryuenc.DecodeFlowEntries(flow), nil
		}
	}

	return nil, fmt.Errorf("Internal error")
}

func (c *RyuClient) GetGroup(dpid int) ([]*ofproto.GroupEntry, error) {
	r, err := HttpGet(fmt.Sprintf("%s/%s/%d", c.url, "groupdesc", dpid))
	if err != nil {
		return nil, err
	}

	groups := make(map[int][]map[string]interface{})
	d := json.NewDecoder(r)
	if err := d.Decode(&groups); err != nil {
		return nil, err
	}

	for id, group := range groups {
		if id == dpid {
			return ryuenc.DecodeGroupEntries(group)
		}
	}

	return nil, fmt.Errorf("Internal error")
}

func (c *RyuClient) ModFlow(cmd string, mod *ofproto.FlowMod) error {
	b, err := json.Marshal(mod)
	if err != nil {
		return err
	}

	return HttpPost(fmt.Sprintf("%s/%s/%s", c.url, "flowentry", cmd), b)
}

func (c *RyuClient) ModGroup(cmd string, mod *ofproto.GroupMod) error {
	b, err := json.Marshal(mod)
	if err != nil {
		return err
	}

	return HttpPost(fmt.Sprintf("%s/%s/%s", c.url, "groupentry", cmd), b)
}
