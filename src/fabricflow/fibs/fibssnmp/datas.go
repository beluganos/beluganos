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
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type StatsDatas interface {
	PortStatsList() PortStatsList
}

//
// DataServer is server for stats data.
//
type DataServer struct {
	psList PortStatsList
	viper  *viper.Viper
}

func (s *DataServer) PortStatsList() PortStatsList {
	return s.psList
}

//
// NewDataServer returns new instance.
//
func NewDataServer() *DataServer {
	return &DataServer{
		psList: PortStatsList{},
		viper:  viper.New(),
	}
}

//
// Update read and replace port stats data.
//
func (d *DataServer) Update() {
	datas := struct {
		PSList PortStatsList `mapstructure:"port_stats"`
	}{}
	if err := d.viper.UnmarshalExact(&datas); err != nil {
		log.Errorf("Unmarshal error. %s", err)
		return
	}

	psList := datas.PSList.Validate()
	psList.Sort()
	d.psList = psList

	for _, ps := range psList {
		log.Debugf("Update PortStats: %v", ps)
	}
}

//
// touchFile create specified file if exists.
//
func (d *DataServer) touchFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			log.Errorf("Create file error. %s", err)
			return err
		}
		defer f.Close()

		log.Debugf("Touch file success. %s", path)
		return nil

	} else if err != nil {
		log.Errorf("Stat file error. %s", err)
		return err
	}

	log.Debugf("File already exist. %s", path)
	return nil
}

//
// Init starts sub modules.
//
func (d *DataServer) Init(path string, format string) error {
	if err := d.touchFile(path); err != nil {
		log.Errorf("Touch error. %s", err)
		return err
	}

	d.viper.SetConfigFile(path)
	d.viper.SetConfigType(format)
	if err := d.viper.ReadInConfig(); err != nil {
		log.Errorf("Read error. %s", err)
		return err
	}

	d.Update()

	d.viper.WatchConfig()
	d.viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("OnConfigChange")
		d.Update()
	})

	log.Debugf("DataServer started. %s %s", path, format)
	return nil
}
