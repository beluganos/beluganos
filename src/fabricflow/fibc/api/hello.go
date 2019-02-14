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
	"github.com/golang/protobuf/proto"
)

//
// Hello
//
func (h *Hello) Type() uint16 {
	return uint16(FFM_HELLO)
}

func (h *Hello) Bytes() ([]byte, error) {
	return proto.Marshal(h)
}

func NewHello(reId string) *Hello {
	return &Hello{
		ReId: reId,
	}
}

func NewHelloFromBytes(data []byte) (*Hello, error) {
	hello := &Hello{}
	if err := proto.Unmarshal(data, hello); err != nil {
		return nil, err
	}

	return hello, nil
}

//
// FFHello
//
func (h *FFHello) Type() uint16 {
	return uint16(FFM_FF_HELLO)
}

func (h *FFHello) Bytes() ([]byte, error) {
	return proto.Marshal(h)
}

func NewFFHello(dpId uint64) *FFHello {
	return &FFHello{
		DpId: dpId,
	}
}

func NewFFHelloFromBytes(data []byte) (*FFHello, error) {
	hello := &FFHello{}
	if err := proto.Unmarshal(data, hello); err != nil {
		return nil, err
	}

	return hello, nil
}
