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

package fflib

import (
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func LXDConnect(f func(lxd.InstanceServer) error) error {
	con, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return err
	}
	defer con.Disconnect()

	return f(con)
}

func LXDContainerList() ([]*api.Container, error) {
	containers := []*api.Container{}

	err := LXDConnect(func(con lxd.InstanceServer) error {
		names, err := con.GetContainerNames()
		if err != nil {
			return err
		}

		for _, name := range names {
			container, _, err := con.GetContainer(name)
			if err != nil {
				return err
			}
			containers = append(containers, container)
		}

		return nil
	})

	return containers, err
}

func LXDContainerGet(name string) (*api.Container, error) {
	var container *api.Container
	err := LXDConnect(func(con lxd.InstanceServer) error {
		var err error
		container, _, err = con.GetContainer(name)
		return err
	})
	return container, err
}

func LXDContainerHostIfnames(name string, excludeDevices []string) ([]string, error) {
	container, err := LXDContainerGet(name)
	if err != nil {
		return nil, err
	}

	ifnames := []string{}

FOR_LOOP:
	for devName, device := range container.ExpandedDevices {
		if IndexOf(devName, excludeDevices) != -1 {
			continue FOR_LOOP
		}

		if nictype, ok := device["nictype"]; !ok || nictype != "p2p" {
			continue FOR_LOOP
		}

		hostName, ok := device["host_name"]
		if !ok {
			continue FOR_LOOP
		}

		ifnames = append(ifnames, hostName)
	}

	return ifnames, nil
}
