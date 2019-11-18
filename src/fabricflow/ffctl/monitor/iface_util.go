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
	log "github.com/sirupsen/logrus"
)

func dumpDataSets(dss *IfaceStatsDataSets) {
	dss.Range(func(dsKey string, ds *IfaceStatsDataSet) {
		log.Debugf("DataSet: key:'%s'  label:'%s'", dsKey, ds.Label)

		ds.Range(func(dataKey string, data *IfaceStatsData) {
			log.Debugf("  Data: key:'%s'  label:'%s'", dataKey, data.Label)

			data.Range(func(valKey string, value *IfaceStatsValue) {
				log.Debugf("    Value: key:'%s'  label:'%s'", valKey, value.Label)
				log.Debugf("      C:%d D:%d A:%f", value.Counter, value.Diff, value.Average)
			})
		})
	})
}
