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

package ofproto

import (
	"fmt"
)

type Desc struct {
	Dp  string `json:"dp_desc"`
	Sw  string `json:"sw_desc"`
	Hw  string `json:"hw_desc"`
	Sn  string `json:"serial_num"`
	Mfr string `json:"mfr_desc"`
}

func (d *Desc) String() string {
	return fmt.Sprintf("Desc(dp=\"%s\", sw=\"%s\", hw=\"%s\", sn=\"%s\", mfr=\"%s\")",
		d.Dp, d.Sw, d.Hw, d.Sn, d.Mfr)
}
