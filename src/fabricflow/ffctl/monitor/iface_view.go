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

package monitor

import (
	"github.com/jroimartin/gocui"
)

//
// Widget is base of Widgets.
//
type Widget struct {
	Name string

	X int
	Y int
	W int
	H int
}

//
// NewWidget returns new Widget.
//
func NewWidget(name string) *Widget {
	return &Widget{
		Name: name,
		X:    -1,
		Y:    -1,
		W:    1,
		H:    1,
	}
}

//
// Move moves (dx, dy) from current position.
//
func (w *Widget) Move(dx, dy int) {
	w.X += dx
	w.Y += dy
}

//
// ModeTo moves to (x, y)
//
func (w *Widget) MoveTo(x, y int) {
	w.X = x
	w.Y = y
}

//
// SetMax resize Widget (X, Y) - (mzxX - X, maxY - Y)
//
func (w *Widget) SetMax(maxX, maxY int) {
	w.Resize(maxX-w.X, maxY-w.Y)
}

//
// Resize resizes widget.
//
func (w *Widget) Resize(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	w.W = width
	w.H = height
}

//
// SetView sets view setting.
//
func (w *Widget) SetView(g *gocui.Gui) (*gocui.View, error) {
	return g.SetView(w.Name, w.X, w.Y, w.X+w.W, w.Y+w.H)
}

//
// IfaceStatsValue is values of statistics.
//
type IfaceStatsValue struct {
	Label   string
	Counter uint64
	Diff    uint64
	Average float64
}

//
// NewIfaceStatsValue returns new IfaceStatsValue.
//
func NewIfaceStatsValue() *IfaceStatsValue {
	return &IfaceStatsValue{}
}

//
// SetDiff sets Diff.
//
func (d *IfaceStatsValue) SetDiff(other *IfaceStatsValue) {
	if d == nil || other == nil {
		return
	}

	if d.Counter < other.Counter {
		d.Diff = 0
	} else {
		d.Diff = d.Counter - other.Counter
	}
}

//
// SetAverage sets average.
//
func (d *IfaceStatsValue) SetAverage(persec float64, others ...*IfaceStatsValue) {
	if d == nil {
		return
	}

	sum := d.Diff
	for _, other := range others {
		sum += other.Diff
	}

	d.Average = float64(sum) / float64(len(others)+1) / persec
}

type IfaceStatsData struct {
	Label  string
	values map[string]*IfaceStatsValue
	keys   []string
}

func NewIfaceStatsData() *IfaceStatsData {
	return &IfaceStatsData{
		values: map[string]*IfaceStatsValue{},
		keys:   []string{},
	}
}

func (d *IfaceStatsData) Value(key string) *IfaceStatsValue {
	v, ok := d.values[key]
	if !ok {
		v = NewIfaceStatsValue()
		d.values[key] = v
		d.keys = append(d.keys, key)
	}
	return v
}

func (d *IfaceStatsData) Range(f func(string, *IfaceStatsValue)) {
	for _, key := range d.keys {
		f(key, d.values[key])
	}
}

func (d *IfaceStatsData) SetDiff(other *IfaceStatsData) {
	if d == nil || other == nil {
		return
	}

	for key, value := range d.values {
		if otherValue, ok := other.values[key]; ok {
			value.SetDiff(otherValue)
		}
	}
}

func (d *IfaceStatsData) SetAverage(persec float64, others ...*IfaceStatsData) {
	if d == nil {
		return
	}

	for key, value := range d.values {
		otherValues := []*IfaceStatsValue{}
		for _, other := range others {
			if otherValue, ok := other.values[key]; ok {
				otherValues = append(otherValues, otherValue)
			}
		}

		value.SetAverage(persec, otherValues...)
	}
}

type IfaceStatsDataSet struct {
	Label string
	datas map[string]*IfaceStatsData
	keys  []string
}

func NewIfaceStatsDataSet() *IfaceStatsDataSet {
	return &IfaceStatsDataSet{
		datas: map[string]*IfaceStatsData{},
		keys:  []string{},
	}
}

func (d *IfaceStatsDataSet) Data(key string) *IfaceStatsData {
	data, ok := d.datas[key]
	if !ok {
		data = NewIfaceStatsData()
		data.Label = key
		d.datas[key] = data
		d.keys = append(d.keys, key)
	}
	return data
}

func (d *IfaceStatsDataSet) Range(f func(string, *IfaceStatsData)) {
	for _, key := range d.keys {
		f(key, d.datas[key])
	}
}

func (d *IfaceStatsDataSet) SetDiff(other *IfaceStatsDataSet) {
	if d == nil || other == nil {
		return
	}

	for key, data := range d.datas {
		if otherData, ok := other.datas[key]; ok {
			data.SetDiff(otherData)
		}
	}
}

func (d *IfaceStatsDataSet) SetAverage(persec float64, others ...*IfaceStatsDataSet) {
	if d == nil {
		return
	}

	for key, data := range d.datas {
		otherDatas := []*IfaceStatsData{}
		for _, other := range others {
			if otherData, ok := other.datas[key]; ok {
				otherDatas = append(otherDatas, otherData)
			}
		}

		data.SetAverage(persec, otherDatas...)
	}
}

type IfaceStatsDataSets struct {
	dataSets map[string]*IfaceStatsDataSet
	keys     []string
}

func NewIfaceStatsDataSets() *IfaceStatsDataSets {
	return &IfaceStatsDataSets{
		dataSets: map[string]*IfaceStatsDataSet{},
		keys:     []string{},
	}
}

func (d *IfaceStatsDataSets) DataSet(key string) *IfaceStatsDataSet {
	ds, ok := d.dataSets[key]
	if !ok {
		ds = NewIfaceStatsDataSet()
		ds.Label = key
		d.dataSets[key] = ds
		d.keys = append(d.keys, key)
	}

	return ds
}

func (d *IfaceStatsDataSets) Select(key string) *IfaceStatsDataSet {
	if ds, ok := d.dataSets[key]; ok {
		return ds
	}
	return nil
}

func (d *IfaceStatsDataSets) Range(f func(string, *IfaceStatsDataSet)) {
	for _, key := range d.keys {
		f(key, d.dataSets[key])
	}
}

func (d *IfaceStatsDataSets) SetDiff(other *IfaceStatsDataSets) {
	if d == nil || other == nil {
		return
	}

	for key, ds := range d.dataSets {
		if otherDs, ok := other.dataSets[key]; ok {
			ds.SetDiff(otherDs)
		}
	}
}

func (d *IfaceStatsDataSets) SetAverage(persec float64, others ...*IfaceStatsDataSets) {
	if d == nil {
		return
	}

	for key, ds := range d.dataSets {
		otherDSs := []*IfaceStatsDataSet{}
		for _, other := range others {
			if otherDS, ok := other.dataSets[key]; ok {
				otherDSs = append(otherDSs, otherDS)
			}
		}

		ds.SetAverage(persec, otherDSs...)
	}
}
