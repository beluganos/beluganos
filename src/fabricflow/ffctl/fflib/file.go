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

package fflib

import (
	"fmt"
	"os"
	"time"
)

func IndexOf(s string, arr []string) int {
	for index, a := range arr {
		if s == a {
			return index
		}
	}

	return -1
}

func CreateFile(path string, overwrite bool, f func(string)) (*os.File, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) && !overwrite {
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().UTC().In(jst).Format("20060102_150405")
		backupPath := fmt.Sprintf("%s_%s", path, now)
		if err := os.Rename(path, backupPath); err != nil {
			return nil, err
		}

		f(backupPath)
	}

	return os.Create(path)
}
