// -*- coding: utf-8 -*-

package dpport

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var EMPTY_TUPLE = errors.New("EMPTY_TUPLE")

var dpPortOpenFlow = map[uint]uint{
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

var dpPortAS5812 = map[uint]uint{
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

var dpPortAS7712x4 = map[uint]uint{
	1: 50, 2: 54, 3: 58, 4: 62,
	5: 68, 6: 72, 7: 76, 8: 80,
	9: 34, 10: 38, 11: 42, 12: 46,
	13: 84, 14: 88, 15: 92, 16: 96,
	17: 102, 18: 106, 19: 110, 20: 114,
	21: 17, 22: 21, 23: 25, 24: 29,
	25: 118, 26: 122, 27: 126, 28: 130,
	29: 1, 30: 5, 31: 9, 32: 13,
}

var dpPortMap = map[string]map[uint]uint{
	"openflow": dpPortOpenFlow,
	"as5812":   dpPortAS5812,
	"as7712x4": dpPortAS7712x4,
}

func PortMap(name string) (map[uint]uint, error) {
	if m, ok := selectPortMap(name); ok {
		return m, nil
	}

	return readPortMap(name)
}

func selectPortMap(name string) (map[uint]uint, bool) {
	m, ok := dpPortMap[name]
	return m, ok
}

func readPortMap(path string) (map[uint]uint, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := struct {
		PortMap   map[uint]uint `yaml:"port_map"`
		PortTuple string        `yaml:"port_tuple"`
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

func parsePortTuple(s string) (uint, uint, error) {
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

	return uint(pport), uint(lport), nil
}

func parsePortTuples(r io.Reader) (map[uint]uint, error) {
	m := map[uint]uint{}

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
