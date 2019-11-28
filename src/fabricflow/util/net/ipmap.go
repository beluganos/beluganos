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

package fflibnet

import (
	"container/list"
	"fmt"
	"net"
	"sync"
)

//
// Generator interface
//
type IPMapGenerator interface {
	New(net.IP) (net.IP, error)
	Free(net.IP, net.IP)
}

//
// IPMap
//
type IPMap struct {
	values    map[string]net.IP
	generator IPMapGenerator

	mutex sync.RWMutex
}

func NewIPMap(generator IPMapGenerator) *IPMap {
	if generator == nil {
		generator = NewIPMapDefaultGenerator()
	}

	return &IPMap{
		values:    make(map[string]net.IP),
		generator: generator,
	}
}

func (m *IPMap) Value(key net.IP) (net.IP, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if v, ok := m.values[key.String()]; ok {
		return v, nil
	}

	v, err := m.generator.New(key)
	if err != nil {
		return nil, err
	}

	m.values[key.String()] = v
	return v, nil
}

func (m *IPMap) Free(key net.IP) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.values[key.String()]
	if !ok {
		return
	}

	delete(m.values, key.String())
	m.generator.Free(key, val)
}

func (m *IPMap) Contains(key net.IP) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.values[key.String()]
	return ok
}

func (m *IPMap) Walk(f func(string, net.IP) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for k, v := range m.values {
		if ok := f(k, v); !ok {
			break
		}
	}
}

//
// Generator (default)
//
type IPMapDefaultGenerator struct{}

func NewIPMapDefaultGenerator() *IPMapDefaultGenerator {
	return &IPMapDefaultGenerator{}
}

func (g *IPMapDefaultGenerator) New(key net.IP) (net.IP, error) {
	return key, nil
}

func (g *IPMapDefaultGenerator) Free(key net.IP, val net.IP) {
	// nothing to do.
}

//
// Generator (tield IP Network)
//
type IPMapIPNetGenerator struct {
	ipgen *IPGenerator
}

func NewIPMapIPNetGenerator(nw *net.IPNet) *IPMapIPNetGenerator {
	return &IPMapIPNetGenerator{
		ipgen: NewIPGenerator(nw),
	}
}

func (g *IPMapIPNetGenerator) New(key net.IP) (net.IP, error) {
	return g.ipgen.NextIP()
}

func (g *IPMapIPNetGenerator) Free(key net.IP, val net.IP) {
	// nothing to do.
}

//
// Generator (Pool)
//
type IPMapPoolGenerator struct {
	pool *list.List
}

func NewIPMapPoolGenerator() *IPMapPoolGenerator {
	return &IPMapPoolGenerator{
		pool: list.New(),
	}
}

func (m *IPMapPoolGenerator) Add(val net.IP) {
	m.pool.PushBack(val)
}

func (m *IPMapPoolGenerator) AddIPNet(nw *net.IPNet) {
	ipgen := NewIPGenerator(nw)
	for {
		ip, err := ipgen.NextIP()
		if err != nil {
			break
		}
		m.Add(ip)
	}
}

func (m *IPMapPoolGenerator) New(key net.IP) (net.IP, error) {
	if m.pool.Len() == 0 {
		return nil, fmt.Errorf("No Pool")
	}

	return m.pool.Remove(m.pool.Front()).(net.IP), nil
}

func (m *IPMapPoolGenerator) Free(key net.IP, val net.IP) {
	m.pool.PushBack(val)
}
