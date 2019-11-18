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

package monitor

import (
	fibcapi "fabricflow/fibc/api"
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type IfaceController struct {
	ConfigFile  string
	ConfigType  string
	ConfigStats string
	ConfigDPath string

	FibcAddr string
	FibcPort uint16

	Interval    time.Duration
	HistorySize uint16

	statsCfg []*IfaceStatsConfig
	DpathCfg *IfaceDPathConfig

	informView *IfaceInformWedget
	statsView  *IfaceStatsWedget
	history    *IfaceStatsHistory

	statsCh chan *IfaceStatsMessage
	done    chan struct{}

	fibcClient fibcapi.FIBCApApiClient

	log *log.Entry
}

func NewIfaceController() *IfaceController {
	return &IfaceController{
		statsCfg: []*IfaceStatsConfig{},
		DpathCfg: NewIfaceDPathConfig(),

		informView: NewIfaceInformWedget(),
		statsView:  NewIfaceStatsWedget(),
		history:    NewIfaceStatsHistory(),

		statsCh: make(chan *IfaceStatsMessage),
		done:    make(chan struct{}),

		log: log.WithFields(log.Fields{"module": "ctrl"}),
	}
}

func (c *IfaceController) getConfig() error {
	cfg := NewIfaceConfig().SetConfig(c.ConfigFile, c.ConfigType)
	if err := cfg.Read(); err != nil {
		return err
	}

	statsCfg, ok := cfg.StatsConfig(c.ConfigStats)
	if !ok {
		return fmt.Errorf("Config not found. stats:%s", c.ConfigStats)
	}

	DpathCfg, ok := cfg.DPathConfig(c.ConfigDPath)
	if !ok {
		return fmt.Errorf("config not found. ports:%s", c.ConfigDPath)
	}

	c.statsCfg = statsCfg

	if c.DpathCfg.DpID == 0 {
		c.DpathCfg.DpID = DpathCfg.DpID
	}

	if c.DpathCfg.Ifaces == nil || len(c.DpathCfg.Ifaces) == 0 {
		c.DpathCfg.Ifaces = DpathCfg.Ifaces
	}

	// DEBUG
	c.log.Debugf("getConfig: dp-id : %d", c.DpathCfg.DpID)
	c.log.Debugf("getConfig: ifaces: %v", c.DpathCfg.Ifaces)
	for _, cfg := range c.statsCfg {
		c.log.Debugf("getConfig: stats: %v", cfg.Label)
		for _, cnt := range cfg.Counters {
			c.log.Debugf("getConfig:   couter: %v", cnt)
		}
	}

	return nil
}

func (c *IfaceController) Run() error {
	if err := c.getConfig(); err != nil {
		return err
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		c.log.Errorf("run: NewGui error. %s", err)
		return err
	}
	defer g.Close()

	title := fmt.Sprintf("dpid:%d  ave:%d  %s",
		c.DpathCfg.DpID,
		c.HistorySize,
		c.Interval,
	)
	c.history.SetInterval(c.Interval)
	c.history.SetHistorySize(c.HistorySize)
	c.informView.SetTitle(title)
	c.statsView.MoveTo(0, 3)
	g.SetManager(c.informView, c.statsView)

	// global Key bindings
	c.setKeyBinding(g)

	if err := c.startPortStats(); err != nil {
		c.log.Errorf("run: c.startPortStats error. %s", err)
		return err
	}
	if err := c.startUpdateGui(g); err != nil {
		c.log.Errorf("run : startUpdateGui error. %s", err)
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		c.log.Errorf("run: MainLoop error. %s", err)
		return err
	}

	return nil
}

func (c *IfaceController) setKeyBinding(g *gocui.Gui) error {
	quit := func(g *gocui.Gui, v *gocui.View) error {
		if c.done != nil {
			close(c.done)
			c.done = nil
		}
		return gocui.ErrQuit
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		c.log.Errorf("run: SetKeyBinding error. %s", err)
		return err
	}

	resetPortStats := func(g *gocui.Gui, v *gocui.View) error {
		return c.resetPortStats()
	}

	if err := g.SetKeybinding("", 'R', gocui.ModNone, resetPortStats); err != nil {
		c.log.Errorf("run: SetKeyBinding error. %s", err)
		return err
	}

	return nil
}

func (c *IfaceController) startUpdateGui(g *gocui.Gui) error {
	go func() {

	FOR_LOOP:
		for {
			select {
			case msg := <-c.statsCh:
				dataSets, err := c.portStatsToDataSets(msg.PortStats)
				if err != nil {
					continue FOR_LOOP
				}

				if msg.MessageOnly {
					g.Update(func(g *gocui.Gui) error {
						c.informView.SetMessage(msg.String())
						return nil
					})
				} else {
					c.history.Add(dataSets)

					g.Update(func(g *gocui.Gui) error {
						c.informView.SetMessage(msg.String())
						c.statsView.SetDataSets(dataSets)
						return nil
					})
				}

			case <-c.done:
				c.log.Debugf("serveUpdateGui: EXIT.")
				break FOR_LOOP
			}
		}
	}()

	return nil
}

func portStatsToMap(portStats []*fibcapi.FFPortStats) map[string]*fibcapi.FFPortStats {
	m := map[string]*fibcapi.FFPortStats{}

	for _, ps := range portStats {
	VALUE_LOOP:
		for k, v := range ps.SValues {
			if k == "ifName" {
				m[v] = ps
				break VALUE_LOOP
			}
		}
	}

	return m
}

func (c *IfaceController) portStatsToDataSets(portStats []*fibcapi.FFPortStats) (*IfaceStatsDataSets, error) {
	dataSets := NewIfaceStatsDataSets()
	psMap := portStatsToMap(portStats)

	for _, ifname := range c.DpathCfg.Ifaces {
		ds := dataSets.DataSet(ifname)
		ps, ok := psMap[ifname]
		if !ok {
			continue
		}

		for _, statsCfg := range c.statsCfg {
			statsData := ds.Data(statsCfg.Label)

			for _, counter := range statsCfg.Counters {
				value := statsData.Value(counter.Name)
				value.Label = counter.Label
				if v, ok := ps.Values[counter.Name]; ok {
					value.Counter = v
				} else {
					value.Counter = 0
				}
			}
		}
	}

	return dataSets, nil
}
