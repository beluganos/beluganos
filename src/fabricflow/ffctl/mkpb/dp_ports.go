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

package mkpb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var EMPTY_TUPLE = errors.New("EMPTY_TUPLE")

var dpPortOpenFlow = map[uint32]uint32{
	1: 1, 2: 2, 3: 3, 4: 4,
	5: 5, 6: 6, 7: 7, 8: 8, 9: 9,
	10: 10, 11: 11, 12: 12, 13: 13, 14: 14,
	15: 15, 16: 16, 17: 17, 18: 18, 19: 19,
	20: 20, 21: 21, 22: 22, 23: 23, 24: 24,
	25: 25, 26: 26, 27: 27, 28: 28, 29: 29,
	30: 30, 31: 31, 32: 32, 33: 33, 34: 34,
	35: 35, 36: 36, 37: 37, 38: 38, 39: 39,
	40: 40, 41: 41, 42: 42, 43: 43, 44: 44,
	45: 45, 46: 46, 47: 47, 48: 48, 49: 49,
	50: 50, 51: 51, 52: 52, 53: 53, 54: 54,
	55: 55, 56: 56, 57: 57, 58: 58, 59: 59,
	60: 60, 61: 61, 62: 62, 63: 63, 64: 64,
}

var dpPortAS5812 = map[uint32]uint32{
	1: 1, 2: 2, 3: 3, 4: 4,
	5: 5, 6: 6, 7: 7, 8: 8, 9: 9,
	10: 10, 11: 11, 12: 12, 13: 13, 14: 14,
	15: 15, 16: 16, 17: 17, 18: 18, 19: 19,
	20: 20, 21: 21, 22: 22, 23: 23, 24: 24,
	25: 25, 26: 26, 27: 27, 28: 28, 29: 29,
	30: 30, 31: 31, 32: 32, 33: 33, 34: 34,
	35: 35, 36: 36, 37: 37, 38: 38, 39: 39,
	40: 40, 41: 41, 42: 42, 43: 43, 44: 44,
	45: 45, 46: 46, 47: 47, 48: 48, 49: 49,
	50: 53, 51: 57, 52: 61, 53: 65, 54: 69,
}

var dpPortAS7712 = map[uint32]uint32{
	1:  50,
	2:  54,
	3:  58,
	4:  62,
	5:  68,
	6:  72,
	7:  76,
	8:  80,
	9:  34,
	10: 38,
	11: 42,
	12: 46,
	13: 84,
	14: 88,
	15: 92,
	16: 96,
	17: 102,
	18: 106,
	19: 110,
	20: 114,
	21: 17,
	22: 21,
	23: 25,
	24: 29,
	25: 118,
	26: 122,
	27: 126,
	28: 130,
	29: 1,
	30: 5,
	31: 9,
	32: 13,
}

var dpPortAS7712x4 = map[uint32]uint32{
	1: 50, 2: 51, 3: 52, 4: 53,
	5: 54, 6: 55, 7: 56, 8: 57,
	9: 58, 10: 59, 11: 60, 12: 61,
	13: 62, 14: 63, 15: 64, 16: 65,
	17: 68, 18: 69, 19: 70, 20: 71,
	21: 72, 22: 73, 23: 74, 24: 75,
	25: 76, 26: 77, 27: 78, 28: 79,
	29: 80, 30: 81, 31: 82, 32: 83,
	33: 34, 34: 35, 35: 36, 36: 37,
	37: 38, 38: 39, 39: 40, 40: 41,
	41: 42, 42: 43, 43: 44, 44: 45,
	45: 46, 46: 47, 47: 48, 48: 49,
	49: 84, 50: 85, 51: 86, 52: 87,
	53: 88, 54: 89, 55: 90, 56: 91,
	57: 92, 58: 93, 59: 94, 60: 95,
	61: 96, 62: 97, 63: 98, 64: 99,
	65: 102, 66: 103, 67: 104, 68: 105,
	69: 106, 70: 107, 71: 108, 72: 109,
	73: 110, 74: 111, 75: 112, 76: 113,
	77: 114, 78: 115, 79: 116, 80: 117,
	81: 17, 82: 18, 83: 19, 84: 20,
	85: 21, 86: 22, 87: 23, 88: 24,
	89: 25, 90: 26, 91: 27, 92: 28,
	93: 29, 94: 30, 95: 31, 96: 32,
	97: 118, 98: 119, 99: 120, 100: 121,
	101: 122, 102: 123, 103: 124, 104: 125,
	105: 126, 106: 127, 107: 128, 108: 129,
	109: 130, 110: 131, 111: 132, 112: 133,
	113: 1, 114: 2, 115: 3, 116: 4,
	117: 5, 118: 6, 119: 7, 120: 8,
	121: 9, 122: 10, 123: 11, 124: 12,
	125: 13, 126: 14, 127: 15, 128: 16,
}

var dpPortMap = map[string]map[uint32]uint32{
	"openflow": dpPortOpenFlow,
	"as5812":   dpPortAS5812,
	"as7712":   dpPortAS7712,
	"as7712x4": dpPortAS7712x4,
}

func PortMapKeys() []string {
	keys := []string{}
	for key := range dpPortMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func PortMap(name string) (map[uint32]uint32, error) {
	if m, ok := selectPortMap(name); ok {
		return m, nil
	}

	return readPortMap(name)
}

func RangePortMap(m map[uint32]uint32, f func(uint32, uint32)) {
	if m == nil {
		return
	}
	for k, v := range m {
		f(k, v)
	}
}

func ConvToPortMap(pports []uint32, portMap map[uint32]uint32) map[uint32]uint32 {
	if len(pports) == 0 {
		return portMap
	}

	m := map[uint32]uint32{}
	for _, pport := range pports {
		if lport, ok := portMap[pport]; ok {
			m[pport] = lport
		}
	}

	return m
}

func ConvToLPortList(pports []uint32, portMap map[uint32]uint32) []uint32 {
	ports := []uint32{}
	if pports == nil || len(pports) == 0 {
		for _, lport := range portMap {
			ports = append(ports, lport)
		}
	} else {
		for _, pport := range pports {
			if lport, ok := portMap[pport]; ok {
				ports = append(ports, lport)
			}
		}
	}

	return ports
}

func ConvToPPortList(pports []uint32, portMap map[uint32]uint32) []uint32 {
	ports := []uint32{}
	if pports == nil || len(pports) == 0 {
		for pport, _ := range portMap {
			ports = append(ports, pport)
		}
	} else {
		for _, pport := range pports {
			if _, ok := portMap[pport]; ok {
				ports = append(ports, pport)
			}
		}
	}

	return ports
}

func HasPPort(pport uint32, portMap map[uint32]uint32) bool {
	_, ok := portMap[pport]
	return ok
}

func selectPortMap(name string) (map[uint32]uint32, bool) {
	m, ok := dpPortMap[name]
	return m, ok
}

func readPortMap(path string) (map[uint32]uint32, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := struct {
		PortMap   map[uint32]uint32 `yaml:"port_map"`
		PortTuple string            `yaml:"port_tuple"`
	}{}
	if err := yaml.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	if len(m.PortMap) == 0 {
		pmap, err := parsePortTuples(strings.NewReader(m.PortTuple))
		if err != nil {
			return nil, err
		}
		m.PortMap = pmap
	}

	dpPortMap[path] = m.PortMap

	return m.PortMap, nil
}

func parsePortTuple(s string) (uint32, uint32, error) {
	line := strings.Trim(s, "{} ")
	line = strings.Replace(line, " ", "", -1)
	if len(line) == 0 {
		return 0, 0, EMPTY_TUPLE
	}

	items := strings.Split(line, ",")

	if n := len(items); n != 2 {
		return 0, 0, fmt.Errorf("Invalid format. '%s'", line)
	}

	pport, err := strconv.ParseUint(items[0], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid pport. '%s'", items[0])
	}

	lport, err := strconv.ParseUint(items[1], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid lport. '%s'", items[1])
	}

	return uint32(pport), uint32(lport), nil
}

func parsePortTuples(r io.Reader) (map[uint32]uint32, error) {
	m := map[uint32]uint32{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		pport, lport, err := parsePortTuple(line)
		if err == EMPTY_TUPLE {
			continue
		}
		if err != nil {
			return nil, err
		}

		m[pport] = lport
	}

	return m, nil
}
