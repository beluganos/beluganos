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

package gonslib

import (
	"github.com/beluganos/go-opennsl/opennsl"

	log "github.com/sirupsen/logrus"
)

//
// PortStatInit initialze opennsl port stats function.
//
func PortStatInit(unit int) error {
	if err := opennsl.StatInit(unit); err != nil {
		log.Errorf("Stat init error. %s", err)
		return err
	}

	log.Infof("Stat init ok.")
	return nil
}

//
// PortStats is list of opennsl.StatVal.
//
type PortStats []opennsl.StatVal

//
// NewPortStats creates new instance.
//
func NewPortStats(names []string) PortStats {
	statVals := PortStats{}
	for _, name := range names {
		statVal, err := opennsl.ParseStatVal(name)
		if err != nil {
			log.Errorf("ParsePortStats error. %s", err)
		} else {
			statVals = append(statVals, statVal)
		}
	}

	return statVals
}

//
// Get gets stats of specified port.
//
func (p PortStats) Get(unit int, port opennsl.Port) (map[string]uint64, error) {
	values, err := opennsl.StatValMultiGet(unit, opennsl.Port(port), p...)
	if err != nil {
		return nil, err
	}

	stats := map[string]uint64{}
	for index, statVal := range p {
		stats[statVal.String()] = values[index]
	}

	return stats, nil
}

//
// GetAll gets Stats of all port.
//
func (p PortStats) GetAll(unit int, ports []opennsl.Port) ([]map[string]uint64, error) {
	statsList := make([]map[string]uint64, len(ports))

	for index, port := range ports {
		stats, err := p.Get(unit, port)
		if err != nil {
			return nil, err
		}
		statsList[index] = stats
	}

	return statsList, nil
}
