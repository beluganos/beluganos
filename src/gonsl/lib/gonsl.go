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
	"github.com/beluganos/go-opennsl/sal"

	log "github.com/sirupsen/logrus"
)

//
// SimInit initialize simlator mode.
//
func SimInit(unit int) {
	// FOR DEBUG ONLY
	log.Infof("Initialize for simulator.")
	if err := opennsl.PortInit(unit); err != nil {
		log.Debugf("opennsl.PortInit error. %s", err)
	}
	if _, err := opennsl.VlanDefaultMustGet(unit).Create(unit); err != nil {
		log.Debugf("opennsl.VlanCreate error. %s", err)
	}
}

//
// DriverInit initialize opennsl driver.
//
func DriverInit(unit int, cfg *ONSLConfig) error {
	var init *sal.Init
	if cfg != nil {
		init = &sal.Init{}
		init.Free()
		init.SetCfgFname(cfg.Config)

		log.Infof("Driver init uses config %s", cfg.Config)
	}

	if err := init.Init(); err != nil {
		log.Errorf("Initializer error. %s", err)
		return err
	}

	if err := PortInit(unit); err != nil {
		log.Errorf("Port Init . %s", err)
		return err
	}

	if err := PortStatInit(unit); err != nil {
		log.Errorf("Stat init error. %s", err)
		return err
	}

	if err := RxInit(unit); err != nil {
		log.Errorf("Rx Init error. %s", err)
		return err
	}

	if err := opennsl.SwitchL3EgressMode.Set(unit, 1); err != nil {
		log.Errorf(".SwitchL3EgressMode.Set %s", err)
		return err
	}

	log.Infof("DriverInit ok. unit=%d", unit)
	return nil
}

//
// DriverExit terminates opennsl driver.
//
func DriverExit() {
	log.Infof("DriverExit()")
	sal.DriverExit()
}
