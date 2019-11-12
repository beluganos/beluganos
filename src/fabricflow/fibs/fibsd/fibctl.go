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

package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	FIBCTypeTCP     = "tcp"
	FIBCTypeGrpc    = "grpc"
	FIBCTypeDefault = FIBCTypeGrpc
)

//
// FIBController is interface of FIBC Client.
//
type FIBController interface {
	Dps() ([]uint64, error)
	PortStats(uint64, []string) (PortStats, error)
}

//
// NewFIBController returns new FIBC Client.
//
func NewFIBController(fibcType string, url string) FIBController {
	switch fibcType {
	case FIBCTypeTCP:
		return NewFIBHttpController(url)
	default:
		return NewFIBGrpcController(url)
	}
}

//
// PortStats is port stats datas.
//
type PortStats []map[string]interface{}

//
// NewPortStats returns new PortStats.
//
func NewPortStats() PortStats {
	return []map[string]interface{}{}
}

func (p PortStats) normalize() {
	for _, ps := range p {
		for key, val := range ps {
			switch v := val.(type) {
			case float64:
				ps[key] = int64(v)
			default:
			}
		}
	}
}

//
// WriteTo writes to writer.
//
func (p PortStats) WriteTo(w io.Writer) (int64, error) {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	data := struct {
		PortStats PortStats `yaml:"port_stats"`
	}{
		PortStats: p,
	}

	if err := encoder.Encode(data); err != nil {
		return 0, err
	}

	return 0, nil
}

func writeToTemp(w io.WriterTo) (string, error) {
	file, err := ioutil.TempFile("", "fibs_stats")
	if err != nil {
		return "", err
	}

	defer file.Close()

	os.Chmod(file.Name(), 0644)

	bfile := bufio.NewWriter(file)
	w.WriteTo(bfile)
	bfile.Flush()

	return file.Name(), nil
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func moveFile(srcPath, dstPath string) error {
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}

	if err := copyFile(srcPath, dstPath); err != nil {
		return err
	}

	return os.Remove(srcPath)
}
