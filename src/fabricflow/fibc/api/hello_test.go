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

package fibcapi

import (
	"testing"
)

func TestHello_Type(t *testing.T) {
	h := Hello{}

	if v := h.Type(); v != uint16(FFM_HELLO) {
		t.Errorf("Hello Type unmatch. %d", v)
	}
}

func TestHello_Bytes(t *testing.T) {
	h := Hello{
		ReId: "1.1.1.1",
	}

	b, err := h.Bytes()

	if err != nil {
		t.Errorf("Hello Bytes error. %s", err)
	}

	if v := len(b); v == 0 {
		t.Errorf("Hello Bytes unmatch. %d", v)
	}
}

func TestHello_New(t *testing.T) {
	h := NewHello("1.1.1.1")

	if v := h.ReId; v != "1.1.1.1" {
		t.Errorf("NewHello ReId unmatch. %s", v)
	}
}
