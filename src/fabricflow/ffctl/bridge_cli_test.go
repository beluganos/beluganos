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
	"fmt"
	"testing"
)

func TestParseBridgeCLI(t *testing.T) {
	data := `[{
        "mac": "33:33:00:00:00:01",
        "dev": "l2swbr0",
        "flags": ["self"
        ],
        "state": "permanent"
    },{
        "mac": "00:16:3e:09:a4:44",
        "dev": "eth1",
        "vlan": 1,
        "master": "l2swbr0",
        "state": "permanent"
    },{
        "mac": "00:16:3e:09:a4:44",
        "dev": "eth1",
        "vlan": 10,
        "master": "l2swbr0",
        "state": "permanent"
    }
]
`
	ret, err := ParseBridgeFdbCLI([]byte(data))
	if err != nil {
		t.Errorf("ParseBridgeCLI error. %s", err)
	}

	if v := len(ret); v != 3 {
		t.Errorf("ParseBridgeCLI unmatch. num=%d", v)
	}

	for _, r := range ret {
		fmt.Printf("%s\n", r)
	}
}
