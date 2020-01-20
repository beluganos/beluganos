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
// Multipart
//
type FFMultipartRequestHandler interface {
	FIBCFFMultipartRequest(*fibcnet.Header, *FFMultipart_Request) error
}

type FFMultipartReplyHandler interface {
	FIBCFFMultipartReply(*fibcnet.Header, *FFMultipart_Reply) error
}

//
// Multipart.Port
//
type FFMultipartPortRequestHandler interface {
	FIBCFFMultipartPortRequest(*fibcnet.Header, *FFMultipart_Request, *FFMultipart_PortRequest)
}

type FFMultipartPortReplyHandler interface {
	FIBCFFMultipartPortReply(*fibcnet.Header, *FFMultipart_Reply, *FFMultipart_PortReply)
}

//
// Multipart.PortDesc
//
type FFMultipartPortDescRequestHandler interface {
	FIBCFFMultipartPortDescRequest(*fibcnet.Header, *FFMultipart_Request, *FFMultipart_PortDescRequest)
}

type FFMultipartPortDescReplyHandler interface {
	FIBCFFMultipartPortDescReply(*fibcnet.Header, *FFMultipart_Reply, *FFMultipart_PortDescReply)
}
