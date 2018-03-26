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

package fibccmd

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func httpJson(url string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return httpPost(url, b)
}

func httpPost(url string, body []byte) error {
	r := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
