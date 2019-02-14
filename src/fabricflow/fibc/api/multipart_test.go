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

package fibcapi

import (
	"testing"
)

func TestNewFFMultipart_Request_Port(t *testing.T) {
	dpId := uint64(0x123456789a)
	portNo := uint32(11)
	req := NewFFMultipart_Request_Port(dpId, portNo)

	if v := req.DpId; v != dpId {
		t.Errorf("NewFFMultipart_Request_Port DpId unmatch. %d/%d", dpId, v)
	}

	if v := req.MpType; v != FFMultipart_PORT {
		t.Errorf("TestNewFFMultipart_Request_Port MpType unmatch. %d", v)
	}

	port := req.GetPort()
	if v := port.PortNo; v != portNo {
		t.Errorf("TestNewFFMultipart_Request_Port PortNo unmatch. %d", v)
	}
}

func TestNewFFMultipart_Reply_Port(t *testing.T) {
	dpId := uint64(0x123456789a)
	stats := []*FFPortStats{
		NewFFPortStats(10, map[string]uint64{}),
	}
	req := NewFFMultipart_Reply_Port(dpId, stats)

	if v := req.DpId; v != dpId {
		t.Errorf("NewFFMultipart_Request_Port DpId unmatch. %d/%d", dpId, v)
	}

	if v := req.MpType; v != FFMultipart_PORT {
		t.Errorf("TestNewFFMultipart_Request_Port MpType unmatch. %d", v)
	}

	port := req.GetPort()
	if v := len(port.Stats); v != 1 {
		t.Errorf("TestNewFFMultipart_Request_Port #Port.Stats unmatch. %d", v)
	}
	if v := port.Stats; v[0] != stats[0] {
		t.Errorf("TestNewFFMultipart_Request_Port Port.Stats unmatch. %v", v)
	}
}
