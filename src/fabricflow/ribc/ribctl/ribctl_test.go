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

package ribctl

import (
	"testing"

	fibcapi "fabricflow/fibc/api"
	"gonla/nlamsg"
)

func TestRIBController_HandlerImpl(t *testing.T) {
	var c interface{} = &RIBController{}

	if _, ok := c.(fibcapi.DpStatusHandler); !ok {
		t.Errorf("RIBController has no handler. (DpStatus)")
	}

	if _, ok := c.(fibcapi.PortStatusHandler); !ok {
		t.Errorf("RIBController has no handler. (PortStatus)")
	}

	if _, ok := c.(nlamsg.NetlinkNodeHandler); !ok {
		t.Errorf("RIBController has no handler. (NetlinkNode)")
	}

	if _, ok := c.(nlamsg.NetlinkLinkHandler); !ok {
		t.Errorf("RIBController has no handler. (NetlinkLnk)")
	}

	if _, ok := c.(nlamsg.NetlinkAddrHandler); !ok {
		t.Errorf("RIBController has no handler. (NetlinkAddr)")
	}

	if _, ok := c.(nlamsg.NetlinkNeighHandler); !ok {
		t.Errorf("RIBController has no handler. (NetlinkNeigh)")
	}

	if _, ok := c.(nlamsg.NetlinkRouteHandler); !ok {
		t.Errorf("RIBController has no handler. (NetlinkRoute)")
	}
}
