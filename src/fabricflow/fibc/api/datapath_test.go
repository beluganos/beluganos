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

//
// DpStatus
//
func TestDpStatus_Type(t *testing.T) {
	d := DpStatus{}

	if v := d.Type(); v != uint16(FFM_DP_STATUS) {
		t.Errorf("DpStatus.Type unmatch %d", v)
	}
}

func TestDpStatus_Bytes(t *testing.T) {
	ds := DpStatus{
		Status: DpStatus_ENTER,
		ReId:   "1.1.1.1",
	}

	b, err := ds.Bytes()

	if err != nil {
		t.Errorf("DpStatus Bytes error. %s", err)
	}

	if v := len(b); v != 11 {
		t.Errorf("DpStatus Bytes unmatch. %v", b)
	}
}

func TestDpStatus_NewFromByte(t *testing.T) {
	ds := &DpStatus{
		Status: DpStatus_ENTER,
		ReId:   "1.1.1.1",
	}
	b, _ := ds.Bytes()

	d, err := NewDpStatusFromBytes(b)

	if err != nil {
		t.Errorf("DpStatus NewDpStatusFromBytes error. %s", err)
	}

	if *ds != *d {
		t.Errorf("DpStatus NewDpStatusFromBytes unmatch. %v", d)
	}
}

func TestDpStatus_NewFromByte_error(t *testing.T) {
	b := []byte{1, 2, 3, 4}

	_, err := NewDpStatusFromBytes(b)

	if err == nil {
		t.Errorf("DpStatus NewDpStatusFromBytes error. %s", err)
	}
}
