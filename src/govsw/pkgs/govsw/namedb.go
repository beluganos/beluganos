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

package govsw

import (
	"regexp"
	"sync"

	log "github.com/sirupsen/logrus"
)

type NameDB struct {
	linkNames  map[string]struct{}
	linkNameRe map[string]*regexp.Regexp
	blackList  map[string]struct{}
	mutex      sync.RWMutex
	log        *log.Entry
}

func NewNameDB() *NameDB {
	return &NameDB{
		linkNames:  make(map[string]struct{}),
		linkNameRe: make(map[string]*regexp.Regexp),
		blackList:  make(map[string]struct{}),
		log:        log.WithField("module", "namedb"),
	}
}

func (db *NameDB) reset() {
	db.linkNames = map[string]struct{}{}
	db.linkNameRe = map[string]*regexp.Regexp{}
	db.blackList = map[string]struct{}{}
}

func (db *NameDB) registerPattern(expr string) {
	re, err := regexp.Compile(expr)
	if err != nil {
		db.log.Warnf("Register Pattern error. '%s' %s", expr, err)
		return
	}

	db.linkNameRe[expr] = re

	db.log.Debugf("Register Pattern: '%s'", expr)
}

func (db *NameDB) RegisterPattern(expr string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.registerPattern(expr)
}

func (db *NameDB) registerIfname(ifname string) {

	db.linkNames[ifname] = struct{}{}

	db.log.Debugf("Register Ifname:'%s'", ifname)
}

func (db *NameDB) RegisterIfname(ifname string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.registerIfname(ifname)
}

func (db *NameDB) registerBlackList(ifname string) {
	db.blackList[ifname] = struct{}{}

	db.log.Debugf("Register BlackList:'%s'", ifname)
}

func (db *NameDB) RegisterBlackList(ifname string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.registerBlackList(ifname)
}

func (db *NameDB) Update(cfg *DpConfig) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.reset()

	for _, ifname := range cfg.Ifaces.Names {
		db.registerIfname(ifname)
	}

	for _, pattern := range cfg.Ifaces.Patterns {
		db.registerPattern(pattern)
	}

	for _, ifname := range cfg.Ifaces.BlackList {
		db.registerBlackList(ifname)
	}

	return nil
}

func (db *NameDB) Delete(s string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	delete(db.linkNameRe, s)
	delete(db.linkNames, s)
	delete(db.blackList, s)
}

func (db *NameDB) Has(ifname string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if _, ok := db.blackList[ifname]; ok {
		return false
	}

	if _, ok := db.linkNames[ifname]; ok {
		return true
	}

	for _, re := range db.linkNameRe {
		if match := re.MatchString(ifname); match {
			return true
		}
	}

	return false
}

func (db *NameDB) Range(f func(kind, name string)) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for name, _ := range db.linkNames {
		f("name", name)
	}

	for expr, _ := range db.linkNameRe {
		f("pattern", expr)
	}

	for name, _ := range db.blackList {
		f("blacklist", name)
	}
}
