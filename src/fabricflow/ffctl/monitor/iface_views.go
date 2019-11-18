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
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
)

const (
	IfaceInformName = "informView"
	IfaceStatsName  = "statsView"
)

//
// IfaceInformWedget displays Informations.
//
type IfaceInformWedget struct {
	*Widget
	title   string
	message string

	mutex sync.Mutex
}

func NewIfaceInformWedget() *IfaceInformWedget {
	return &IfaceInformWedget{
		Widget:  NewWidget(IfaceInformName),
		message: "plawse wait....",
	}
}

func (w *IfaceInformWedget) SetTitle(title string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.title = title
}

func (w *IfaceInformWedget) SetMessage(msg string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.message = msg
}

func (w *IfaceInformWedget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	w.MoveTo(0, 0)
	w.SetMax(maxX-1, 3)

	v, err := w.SetView(g)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, w.title)
		return nil
	}

	v.Clear()

	w.mutex.Lock()
	defer w.mutex.Unlock()

	fmt.Fprintf(v, " %s\n", w.title)
	fmt.Fprintf(v, " %s\n", w.message)

	return nil
}

//
// IfaceStatsWedget displays stats data.
//
type IfaceStatsWedget struct {
	*Widget
	ds *IfaceStatsDataSets

	mutex sync.Mutex
}

func NewIfaceStatsWedget() *IfaceStatsWedget {
	return &IfaceStatsWedget{
		Widget: NewWidget(IfaceStatsName),
		ds:     NewIfaceStatsDataSets(),
	}
}

func (w *IfaceStatsWedget) SetDataSets(ds *IfaceStatsDataSets) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ds = ds
}

func (w *IfaceStatsWedget) Layout(g *gocui.Gui) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.ds == nil {
		return nil
	}

	maxX, maxY := g.Size()
	w.SetMax(maxX-1, maxY-1)

	v, err := w.SetView(g)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, "plawse wait....")
		return nil
	}

	v.Clear()

	fmt.Fprintln(v, "                                         count    diff          ave")
	w.ds.Range(func(dsKey string, dataSet *IfaceStatsDataSet) {
		fmt.Fprintf(v, "%s\n", dataSet.Label)

		dataSet.Range(func(dataKey string, data *IfaceStatsData) {
			fmt.Fprintf(v, "  %s\n", data.Label)

			data.Range(func(name string, value *IfaceStatsValue) {
				fmt.Fprintf(v, "    %-32s %8d %8d %12.4f\n",
					value.Label,
					value.Counter,
					value.Diff,
					value.Average,
				)
			})
		})
	})

	return nil
}
