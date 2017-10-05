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
//  DpStatus
//
func (*DpStatus) Type() uint16 {
	return uint16(FFM_DP_STATUS)
}

func (d *DpStatus) Bytes() ([]byte, error) {
	return proto.Marshal(d)
}

func NewDpStatusFromBytes(data []byte) (*DpStatus, error) {
	ds := &DpStatus{}
	if err := proto.Unmarshal(data, ds); err != nil {
		return nil, err
	}

	return ds, nil
}
