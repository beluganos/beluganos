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

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

//
// DpsMsg dpid http message.
//
type DpsMsg struct {
	DpIds []uint64 `json:"dpids"`
}

//
// HTTPGet gets dpid message from fibcd.
//
func (m *DpsMsg) HTTPGet(baseURL string) error {
	url := fmt.Sprintf("%s/fib/dps", baseURL)

	res, err := HTTPGet(url)
	if err != nil {
		log.Errorf("HttpGet error. %s %s", url, err)
		return err
	}

	decoder := json.NewDecoder(res)
	if err := decoder.Decode(m); err != nil {
		log.Errorf("json.Decode error. %s", err)
		return err
	}

	return nil
}

//
// PortStats is port stats datas.
//
type PortStats []map[string]interface{}

//
// PortStatsMsg is port stats message.
//
type PortStatsMsg map[string]PortStats

//
// HTTPGet gets port stats by http from fibcd.
//
func (m PortStatsMsg) HTTPGet(baseURL string, dpid uint64) error {
	url := fmt.Sprintf("%s/fib/stats/port/%d", baseURL, dpid)

	res, err := HTTPGet(url)
	if err != nil {
		log.Errorf("HttpGet error. %s %s", url, err)
		return err
	}

	decoder := json.NewDecoder(res)
	if err := decoder.Decode(&m); err != nil {
		log.Errorf("json.Decode error. %s", err)
		return err
	}

	return nil
}

//
// PortStats returns port stats.
//
func (m PortStatsMsg) PortStats(dpid uint64) (PortStats, bool) {
	key := fmt.Sprintf("%d", dpid)
	if ps, ok := m[key]; ok {
		return ps, true
	}

	return nil, false
}

//
// WriteTo writes to writer.
//
func (p PortStats) WriteTo(w io.Writer) (int64, error) {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	data := struct {
		PortStats PortStats `yaml:"port_stats"`
	}{
		PortStats: p,
	}

	if err := encoder.Encode(data); err != nil {
		return 0, err
	}

	return 0, nil
}
