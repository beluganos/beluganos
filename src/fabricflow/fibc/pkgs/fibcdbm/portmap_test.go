// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

package fibcdbm

import "testing"

func TestPortMapListByParent(t *testing.T) {
	m := NewPortMap()

	eth1 := &PortEntry{
		Key: NewPortKey("1.1.1.1", "eth1"),
	}

	eth2 := &PortEntry{
		Key: NewPortKey("1.1.1.1", "eth2"),
	}

	eth2_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth2.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth2"),
	}

	eth3 := &PortEntry{
		Key: NewPortKey("1.1.1.1", "eth3"),
	}

	eth3_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth3.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth3"),
	}

	eth3_20 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth3.20"),
		ParentKey: NewPortKey("1.1.1.1", "eth3"),
	}

	m.Register(eth1)
	m.Register(eth2)
	m.Register(eth2_10)
	m.Register(eth3)
	m.Register(eth3_10)
	m.Register(eth3_20)

	var cnt int

	// eth1
	cnt = 0
	m.ListByParent(NewPortKey("1.1.1.1", "eth1"), func(e *PortEntry) bool {
		cnt++
		return true
	})
	if cnt != 0 {
		t.Errorf("ListByParent unmath. cot=%d", cnt)
	}

	// eth2
	cnt = 0
	m.ListByParent(NewPortKey("1.1.1.1", "eth2"), func(e *PortEntry) bool {
		cnt++
		return true
	})
	if cnt != 1 {
		t.Errorf("ListByParent unmath. cnt%d", cnt)
	}

	// eth3
	cnt = 0
	m.ListByParent(NewPortKey("1.1.1.1", "eth3"), func(e *PortEntry) bool {
		cnt++
		return true
	})
	if cnt != 2 {
		t.Errorf("ListByParent unmath. cnt=%d", cnt)
	}
}

func TestPortMapListByParentMulti(t *testing.T) {
	m := NewPortMap()

	eth1 := &PortEntry{
		Key: NewPortKey("1.1.1.1", "eth1"),
	}

	eth1_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth1.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth1"),
	}

	eth1_10_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth1.10.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth1.10"),
	}

	eth2 := &PortEntry{
		Key: NewPortKey("1.1.1.1", "eth2"),
	}

	eth2_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth2.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth2"),
	}

	eth2_10_10 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth2.10.10"),
		ParentKey: NewPortKey("1.1.1.1", "eth2.10"),
	}

	eth2_20 := &PortEntry{
		Key:       NewPortKey("1.1.1.1", "eth2.20"),
		ParentKey: NewPortKey("1.1.1.1", "eth2"),
	}

	m.Register(eth1)
	m.Register(eth1_10)
	m.Register(eth1_10_10)
	m.Register(eth2)
	m.Register(eth2_10)
	m.Register(eth2_10_10)
	m.Register(eth2_20)

	var cnt int

	// eth1
	cnt = 0
	m.ListByParent(NewPortKey("1.1.1.1", "eth1"), func(e *PortEntry) bool {
		cnt++
		return true
	})
	if cnt != 2 {
		t.Errorf("ListByParent unmath. cnt=%d", cnt)
	}

	// eth2
	cnt = 0
	m.ListByParent(NewPortKey("1.1.1.1", "eth2"), func(e *PortEntry) bool {
		cnt++
		return true
	})
	if cnt != 3 {
		t.Errorf("ListByParent unmath. cnt=%d", cnt)
	}
}
