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

import fibcnet "fabricflow/fibc/net"

//
// OAM(Request)
//
type OAMRequestHandler interface {
	FIBCOAMRequest(*fibcnet.Header, *OAM_Request) error
}

//
// OAM(Reply)
//
type OAMReplyHandler interface {
	FIBCOAMReply(*fibcnet.Header, *OAM_Reply) error
}

//
// OAM.AuditRouteCnt(Request)
//
type OAMAuditRouteCntRequestHandler interface {
	FIBCOAMAuditRouteCntRequest(*fibcnet.Header, *OAM_Request, *OAM_AuditRouteCntRequest) error
}

//
// OAM.AuditRouteCnt(Reply)
//
type OAMAuditRouteCntReplyHandler interface {
	FIBCOAMAuditRouteCntReply(*fibcnet.Header, *OAM_Reply, *OAM_AuditRouteCntReply) error
}
