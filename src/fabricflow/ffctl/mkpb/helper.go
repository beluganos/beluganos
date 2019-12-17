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

package mkpb

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func createFile(path string, overwrite bool, f func(string)) (*os.File, error) {
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

func SortUint32s(v []uint32) {
	sort.Slice(v, func(i, j int) bool {
		return v[i] < v[j]
	})
}

func StrUint32Slice(vv []uint32) string {
	ss := []string{}
	for _, v := range vv {
		ss = append(ss, fmt.Sprintf("%d", v))
	}

	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func StrUint16Slice(vv []uint16) string {
	ss := []string{}
	for _, v := range vv {
		ss = append(ss, fmt.Sprintf("%d", v))
	}

	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func StrStringSlice(ss []string) string {
	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func StrOid(oid []uint32) string {
	ss := make([]string, len(oid)+1)
	for index, n := range oid {
		ss[index+1] = fmt.Sprintf("%d", n)
	}
	return strings.Join(ss, ".")
}
