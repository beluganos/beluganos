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
	"fabricflow/fibc/cmd/lib"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Arg map[string]interface{}

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Read error. %s", err)
		return
	}

	var entries []*fibccmd.ModEntry
	if err := yaml.Unmarshal(buf, &entries); err != nil {
		fmt.Printf("Unmarshal error. %s\n", err)
		return
	}

	c := fibccmd.NewModClient("127.0.0.1:8080")
	if err := c.Sends(entries); err != nil {
		fmt.Printf("Send error. %s\n", err)
	}
}
