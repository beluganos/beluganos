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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type Arg map[string]interface{}

type Entry struct {
	Name string                   `yaml:"name" json:"name"`
	Cmd  string                   `yaml:"cmd" json:"cmd"`
	ReId string                   `yaml:"re_id" json:"re_id"`
	VsId uint64                   `yaml:"vs_id" json:"vs_id"`
	DpId uint64                   `yaml:"dp_id" json:"dp_id"`
	Args []map[string]interface{} `yaml:"args" json:"args"`
}

func Send(addr string, entry *Entry) error {
	url := fmt.Sprintf("http://%s/fib/%s", addr, entry.Name)
	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return HttpPost(url, body)
}

func HttpPost(url string, body []byte) error {
	sr := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", sr)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Read error. %s", err)
		return
	}

	var entries []*Entry
	if err := yaml.Unmarshal(buf, &entries); err != nil {
		fmt.Printf("Unmarshal error. %s\n", err)
		return
	}

	for _, entry := range entries {
		if err := Send("127.0.0.1:8080", entry); err != nil {
			fmt.Printf("Send error. %s\n", err)
		}
	}
}
