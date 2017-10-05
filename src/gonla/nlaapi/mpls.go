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

package nlaapi

import (
	"gonla/nladbm"
)

//
// Mpls (Key)
//

func (k *MplsKey) ToNative() *nladbm.MplsKey {
	return &nladbm.MplsKey{
		NId:    uint8(k.NId),
		LLabel: k.LLabel,
	}
}

func NewMplsKeyFromNative(n *nladbm.MplsKey) *MplsKey {
	return &MplsKey{
		NId:    uint32(n.NId),
		LLabel: n.LLabel,
	}
}
