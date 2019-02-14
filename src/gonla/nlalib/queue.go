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

package nlalib

import (
	"container/list"
	"fmt"
	"sync"
)

const (
	QUEUE_INVALID = -1
)

type Queue struct {
	mutex   sync.Mutex
	cond    *sync.Cond
	queue   *list.List
	maxSize int
}

func NewQueue(maxSize int) *Queue {
	q := &Queue{
		queue:   list.New(),
		maxSize: maxSize,
	}
	q.cond = sync.NewCond(&q.mutex)
	return q
}

func (q *Queue) size() int {
	return q.queue.Len()
}

func (q *Queue) isActive() bool {
	return (q.maxSize != QUEUE_INVALID)
}

func (q *Queue) isEmpty() bool {
	return (q.size() == 0)
}

func (q *Queue) Close() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.maxSize = QUEUE_INVALID
	q.cond.Signal()
}

func (q *Queue) Push(v interface{}) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.isActive() {
		return fmt.Errorf("Queue closed.")
	}

	if q.maxSize > 0 {
		if size := q.size(); size >= q.maxSize {
			return fmt.Errorf("Queue full. %d/%d", size, q.maxSize)
		}
	}

	q.queue.PushBack(v)
	q.cond.Signal()

	return nil
}

func (q *Queue) Pop() (v interface{}, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.isActive() && q.isEmpty() {
		q.cond.Wait()
	}

	if !q.isActive() {
		err = fmt.Errorf("Queue closed.")
		return
	}

	v = q.queue.Remove(q.queue.Front())
	return
}
