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

package mqueue

import (
	"testing"
)

type TestMultiQueueItem struct {
	Name  string
	Value string
}

func (i *TestMultiQueueItem) MQueName() string {
	return i.Name
}

func NewTestMultiQueueItem(name, value string) *TestMultiQueueItem {
	return &TestMultiQueueItem{
		Name:  name,
		Value: value,
	}
}

func TestMultiQueue_push_pop(t *testing.T) {

	var item MultiQueueItem
	itemA1 := NewTestMultiQueueItem("A", "a1")

	ques := NewMultiQueues()
	ques.Push(itemA1)

	if n := ques.Size("A"); n != 1 {
		t.Errorf("Queues.Length unmatch. %d", n)
	}

	item = ques.Pop("A")
	if item == nil {
		t.Errorf("Queues.Pop error. %v", item)
	}
	if v := item.(*TestMultiQueueItem).Value; v != "a1" {
		t.Errorf("Queues.Pop unmatch. value=%v", v)
	}

	if l := ques.Len(); l != 0 {
		t.Errorf("Queues.Pop unmatch. len=%d", l)
	}

	item = ques.Pop("B")
	if item != nil {
		t.Errorf("Queues.Pop error. %v", item)
	}
}

func TestMultiQueue_remove(t *testing.T) {

	var cnt int
	var item MultiQueueItem
	itemA1 := NewTestMultiQueueItem("A", "a1")
	itemA2 := NewTestMultiQueueItem("A", "a2")
	itemA3 := NewTestMultiQueueItem("A", "a3")
	itemB1 := NewTestMultiQueueItem("B", "b1")
	itemB2 := NewTestMultiQueueItem("B", "b2")
	itemC1 := NewTestMultiQueueItem("C", "c1")
	items := []MultiQueueItem{
		itemA1, itemA2, itemA3,
		itemB1, itemB2, itemC1,
	}

	ques := NewMultiQueues()
	for _, item = range items {
		ques.Push(item)
	}

	if l := ques.Len(); l != 3 {
		t.Errorf("Queues.Len unmatch. %d", l)
	}
	if l := ques.Size("A"); l != 3 {
		t.Errorf("Queues.Len unmatch. #A=%d", l)
	}
	if l := ques.Size("B"); l != 2 {
		t.Errorf("Queues.Len unmatch. #B=%d", l)
	}
	if l := ques.Size("C"); l != 1 {
		t.Errorf("Queues.Len unmatch. #C=%d", l)
	}

	cnt = 0
	ques.RemoveByName("A", func(item MultiQueueItem) bool {
		cnt += 1
		return false
	})
	if cnt != 3 {
		t.Errorf("ques.RemoveByName unmatch. cnt=%d", cnt)
	}
	if l := ques.Len(); l != 3 {
		t.Errorf("Queues.Len unmatch. %d", l)
	}
	if l := ques.Size("A"); l != 3 {
		t.Errorf("Queues.Len unmatch. #A=%d", l)
	}

	cnt = 0
	ques.RemoveByName("A", func(item MultiQueueItem) bool {
		cnt += 1
		return true
	})
	if cnt != 3 {
		t.Errorf("ques.RemoveByName unmatch. cnt=%d", cnt)
	}
	if l := ques.Len(); l != 2 {
		t.Errorf("Queues.Len unmatch. %d", l)
	}
	if l := ques.Size("A"); l != 0 {
		t.Errorf("Queues.Len unmatch. #A=%d", l)
	}
}
