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

package fibcsrv

import (
	"fabricflow/fibc/pkgs/fibccfg"
	"fabricflow/fibc/pkgs/fibcdbm"
)

//
// RegisterConfig add config data to db.
//
func (c *DBCtl) RegisterConfig(cfg *fibccfg.Config) {
	c.IDMap().VerUp()
	c.PortMap().VerUp()

	c.registerIDMap(cfg)
	c.registerPortMap(cfg)

	c.PortMap().GC(func(e *fibcdbm.PortEntry) bool {
		c.log.Debugf("PortMap: GC %s", e.Key)
		return false
	})
	c.IDMap().GC(func(e *fibcdbm.IDEntry) bool {
		c.log.Debugf("IDMap: GC %s", e)
		return false
	})
}

func (c *DBCtl) registerIDMap(cfg *fibccfg.Config) {
	for _, router := range cfg.Routers {
		dpcfg, ok := cfg.DPathConfig(router.DPath)
		if !ok {
			c.log.Warnf("DP not found. %s", router.DPath)
			continue
		}

		e := fibcdbm.NewIDEntry(dpcfg.DpID, router.ReID)
		c.IDMap().SelectOrRegister(e, func(old *fibcdbm.IDEntry) bool {
			c.log.Debugf("IDMap: %s", e)
			return true
		})
	}
}

func (c *DBCtl) registerPortMap(cfg *fibccfg.Config) {
	for _, router := range cfg.Routers {
		dpID, ok := c.IDMap().SelectByReID(router.ReID)
		if !ok {
			c.log.Warnf("DP not found. %s", router.DPath)
			continue
		}

		for _, port := range router.Ports {
			key := fibcdbm.NewPortKey(router.ReID, port.Name)
			vmport := fibcdbm.NewPortValueR(router.ReID, 0, false)
			dpport := fibcdbm.NewPortValue(dpID, uint32(port.PortID), false)
			vsport := fibcdbm.NewPortValueR("", 0, false)

			e := &fibcdbm.PortEntry{
				Key:       key,
				ParentKey: nil,
				MasterKey: nil,

				VMPort: vmport,
				DPPort: dpport,
				VSPort: vsport,
			}

			c.PortMap().SelectOrRegister(e, func(old *fibcdbm.PortEntry) bool {
				c.log.Debugf("PortMap: %s", old.Key)
				return true
			})
		}
	}
}
