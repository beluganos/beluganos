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

package main

//
// StatsHandler is interface of StatsHandlers
//
type StatsHandler interface {
	Oid() string
	Get(string) *SnmpReply
	GetNext(string) *SnmpReply
}

//
// NewStatsHandler returns handler instance.
//
func NewStatsHandler(cfg *HandlerConfig, datas StatsDatas) StatsHandler {
	return NewPortStatsHandlerFromConfig(cfg, datas)
}

//
// NewStatssHandlers returns handler list.
//
func NewStatsHandlers(cfgs []*HandlerConfig, datas StatsDatas) ([]StatsHandler, error) {
	handlers := make([]StatsHandler, len(cfgs))
	for index, cfg := range cfgs {
		handlers[index] = NewStatsHandler(cfg, datas)
	}
	return handlers, nil
}
