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

package fibslib

import (
	"testing"
)

func TestReplaceOID_short_oid(t *testing.T) {
	orgOid := []uint{1, 2, 3}
	preOid := []uint{1, 2}
	newOid := []uint{5, 6, 7}

	oid, ok := ReplaceOID(orgOid, preOid, newOid)

	if ok {
		t.Errorf("ReplaceOID must not ok")
	}

	if d := CompareOID(oid, []uint{5, 6, 7, 3}); d != 0 {
		t.Errorf("ReplaceOID not match.")
	}
}

func TestReplaceOID_same_oid(t *testing.T) {
	orgOid := []uint{1, 2, 3}
	preOid := []uint{1, 2, 3}
	newOid := []uint{4, 5, 6}

	oid, ok := ReplaceOID(orgOid, preOid, newOid)

	if !ok {
		t.Errorf("ReplaceOID error.")
	}

	if d := CompareOID(oid, []uint{4, 5, 6}); d != 0 {
		t.Errorf("ReplaceOID not match.")
	}
}

func TestReplaceOID_long_oid(t *testing.T) {
	orgOid := []uint{1, 2, 3, 7}
	preOid := []uint{1, 2, 3}
	newOid := []uint{4, 5, 6}

	oid, ok := ReplaceOID(orgOid, preOid, newOid)

	if !ok {
		t.Errorf("ReplaceOID error.")
	}

	if d := CompareOID(oid, []uint{4, 5, 6, 7}); d != 0 {
		t.Errorf("ReplaceOID not match.")
	}
}
