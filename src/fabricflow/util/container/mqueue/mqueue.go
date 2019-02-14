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
	"container/list"
	"sync"
)

//
// MultiQueueItem
//
type MultiQueueItem interface {
	MQueName() string
}

//
// MultiQueuea
//
type MultiQueues struct {
	mutex  sync.Mutex
	queues map[string]*list.List
}

func NewMultiQueues() *MultiQueues {
	return &MultiQueues{
		queues: map[string]*list.List{},
	}
}

func (mq *MultiQueues) queue(name string) (que *list.List, ok bool) {
	que, ok = mq.queues[name]
	return
}

func (mq *MultiQueues) Size(name string) int {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if que, ok := mq.queues[name]; ok {
		return que.Len()
	}
	return 0
}

func (mq *MultiQueues) Len() int {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	return len(mq.queues)
}

func (mq *MultiQueues) push(item MultiQueueItem) {
	que, ok := mq.queue(item.MQueName())
	if !ok {
		que = list.New()
		mq.queues[item.MQueName()] = que
	}

	que.PushBack(item)
}

func (mq *MultiQueues) Push(item MultiQueueItem) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	mq.push(item)
}

func (mq *MultiQueues) pop(name string) (item MultiQueueItem) {
	if que, ok := mq.queue(name); ok {
		if elm := que.Front(); elm != nil {
			item = elm.Value.(MultiQueueItem)
			que.Remove(elm)
		}

		if que.Len() == 0 {
			delete(mq.queues, name)
		}
	}

	return
}

func (mq *MultiQueues) Pop(name string) MultiQueueItem {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	return mq.pop(name)
}

func (mq *MultiQueues) removeQueue(name string, que *list.List, f func(MultiQueueItem) bool) {
	elm := que.Front()
	for elm != nil {
		next := elm.Next()

		if ok := f(elm.Value.(MultiQueueItem)); ok {
			que.Remove(elm)
		}

		elm = next
	}

	if que.Len() == 0 {
		delete(mq.queues, name)
	}
}

func (mq *MultiQueues) removeByName(name string, f func(MultiQueueItem) bool) {
	if que, ok := mq.queue(name); ok {
		mq.removeQueue(name, que, f)
	}
}

func (mq *MultiQueues) RemoveByName(name string, f func(MultiQueueItem) bool) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	mq.removeByName(name, f)
}

func (mq *MultiQueues) Remove(f func(MultiQueueItem) bool) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for name, que := range mq.queues {
		mq.removeQueue(name, que, f)
	}
}

func (mq *MultiQueues) walkQueue(que *list.List, f func(MultiQueueItem)) {
	elm := que.Front()
	for elm != nil {
		f(elm.Value.(MultiQueueItem))
		elm = elm.Next()
	}
}

func (mq *MultiQueues) WalkByName(name string, f func(MultiQueueItem)) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if que, ok := mq.queue(name); ok {
		mq.walkQueue(que, f)
	}
}

func (mq *MultiQueues) Walk(f func(string, MultiQueueItem)) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for name, que := range mq.queues {
		mq.walkQueue(que, func(item MultiQueueItem) {
			f(name, item)
		})
	}
}

func (mq *MultiQueues) FilterWalk(f func(string, MultiQueueItem) bool) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for name, que := range mq.queues {
		if ok := f(name, nil); ok {
			mq.removeQueue(name, que, func(item MultiQueueItem) bool {
				return f(name, item)
			})
		}
	}
}
