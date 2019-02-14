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

type Args struct {
	NSInterval time.Duration
	Verbose    bool
}

func NewArgs() (*Args, error) {
	args := &Args{}
	if err := args.Parse(); err != nil {
		return nil, err
	}

	return args, nil
}

func (a *Args) Parse() error {
	flag.DurationVarP(&a.NSInterval, "ns-interval", "", 15*time.Minute, "sending NS interval.")
	flag.BoolVarP(&a.Verbose, "verbose", "v", false, "show detail messages.")
	flag.Parse()
	return nil
}

func main() {

	args, err := NewArgs()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Hello")

	done := make(chan struct{})
	s := NewServer()
	s.SetNSInterval(args.NSInterval)
	if err := s.Start(done); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	<-done
}
