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

package gonslib

import (
	"fabricflow/fibc/lib"
	"testing"
)

func TestFlowHandlers(t *testing.T) {
	var s interface{} = &Server{}

	// VLANFlow
	if _, ok := s.(fibclib.FIBCVLANFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCVLANFlowModHandler)")
	}

	// TerminationMacFlow
	if _, ok := s.(fibclib.FIBCTerminationMacFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCTerminationMacFlowModHandler)")
	}

	// MPLSFlow
	if _, ok := s.(fibclib.FIBCMPLSFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSFlowModHandler)")
	}

	// UnicastRoutingFlow
	if _, ok := s.(fibclib.FIBCUnicastRoutingFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCUnicastRoutingFlowModHandler)")
	}

	// BridgingFlow
	if _, ok := s.(fibclib.FIBCBridgingFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCBridgingFlowModHandler)")
	}

	// PolicyACLFlow
	if _, ok := s.(fibclib.FIBCPolicyACLFlowModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCPolicyACLFlowModHandler)")
	}
}

func TestGroupHandlers(t *testing.T) {
	var s interface{} = &Server{}

	// L2InterfaceGroup
	if _, ok := s.(fibclib.FIBCL2InterfaceGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCL2InterfaceGroupModHandler)")
	}

	// L3UnicastGroup
	if _, ok := s.(fibclib.FIBCL3UnicastGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCL3UnicastGroupModHandler)")
	}

	// MPLSInterfaceGroup
	if _, ok := s.(fibclib.FIBCMPLSInterfaceGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSInterfaceGroupModHandler)")
	}

	// MPLSLabelGroup
	if _, ok := s.(fibclib.FIBCMPLSLabelL2VpnGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSLabelL2VpnGroupModHandler)")
	}

	if _, ok := s.(fibclib.FIBCMPLSLabelL3VpnGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSLabelL3VpnGroupModHandler)")
	}

	if _, ok := s.(fibclib.FIBCMPLSLabelTun1GroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSLabelTun1GroupModHandler)")
	}

	if _, ok := s.(fibclib.FIBCMPLSLabelTun2GroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSLabelTun2GroupModHandler)")
	}

	if _, ok := s.(fibclib.FIBCMPLSLabelSwapGroupModHandler); !ok {
		t.Errorf("Server not implement handler(FIBCMPLSLabelSwapGroupModHandler)")
	}
}

func TestHandlers(t *testing.T) {
	var s interface{} = &Server{}

	if _, ok := s.(fibclib.FFPacketOutHandler); !ok {
		t.Errorf("erver not implement handler(FFPacketOutHandler)")
	}

	if _, ok := s.(fibclib.FFMultipartPortRequestHandler); !ok {
		t.Errorf("erver not implement handler(FFMultipartPortRequestHandler)")
	}

	if _, ok := s.(fibclib.FFMultipartPortDescRequestHandler); !ok {
		t.Errorf("erver not implement handler(FFMultipartPortDescRequestHandler)")
	}
}
