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
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	appConfigPath = "/etc/beluganos/ribnd.yaml"
	appConfigType = "yaml"
	appConfigName = "default"
)

type App struct {
	ConfigPath string
	ConfigType string
	ConfigName string

	NSInterval time.Duration

	Verbose bool
	Trace   bool

	log *log.Entry
}

func newApp() *App {
	app := &App{
		log: log.WithFields(log.Fields{"module": "app"}),
	}
	app.parse()
	return app
}

func (a *App) parse() {
	flag.StringVarP(&a.ConfigPath, "config-path", "c", appConfigPath, "config file path.")
	flag.StringVarP(&a.ConfigType, "config-type", "", appConfigType, "config file type.")
	flag.StringVarP(&a.ConfigName, "config-name", "", appConfigName, "config name.")
	flag.DurationVarP(&a.NSInterval, "ns-interval", "", 15*time.Minute, "sending NS interval.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")
	flag.BoolVarP(&a.Trace, "trace", "", false, "show more detail messages.")

	flag.Parse()
}

func (a *App) run() error {
	if a.Trace {
		log.SetLevel(log.TraceLevel)
	} else if a.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	a.log.Infof("START")

	done := make(chan struct{})
	defer close(done)

	s := NewServer()
	s.SetNSInterval(a.NSInterval)
	if err := s.SetConfig(a.ConfigPath, a.ConfigType, a.ConfigName); err != nil {
		a.log.Errorf("Config read error. %s", err)
		return err
	}

	if err := s.Start(done); err != nil {
		log.Errorf("Server start error. %s", err)
		return err
	}

	<-done

	return nil
}

func main() {
	if err := newApp().run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
