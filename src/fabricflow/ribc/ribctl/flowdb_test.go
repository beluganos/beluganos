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

package ribctl

import (
	"fmt"
	"testing"
)

func TestFlowDBLoad(t *testing.T) {
	db := NewFlowDB()
	if err := db.SetConfigFile("./flowdb_test.yml", "yaml").Load(); err != nil {
		t.Errorf("Load error. %s", err)
	}

	fmt.Printf("%v\n", db)
	for name, cfg := range db.Flows {
		fmt.Printf("%s %v\n", name, cfg)
		for _, acl := range cfg.PolicyACL {
			fmt.Printf("ACL %v\n", acl)
		}
	}
}

func TestBuiltinFlowConfig(t *testing.T) {
	cfg := NewBuiltinFlowConfig()

	for _, c := range cfg.PolicyACL {
		acl := c.ToAPI()
		if acl.Match == nil {
			t.Errorf("PolicyACL ToAPI match error. %v", c)
		}
		if acl.Action == nil {
			t.Errorf("PolicyACL ToAPI action error. %v", c)
		}
		// fmt.Printf("%v\n", acl)
	}
}
