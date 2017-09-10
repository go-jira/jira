// Copyright (c) 2016 Matt Ho <matt.ho@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jq_test

import (
	"testing"

	"github.com/savaki/jq"
)

func TestParse(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Op       string
		Expected string
		HasError bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Op:       ".hello",
			Expected: `"world"`,
		},
		"nested": {
			In:       `{"a":{"b":"world"}}`,
			Op:       ".a.b",
			Expected: `"world"`,
		},
		"index": {
			In:       `["a","b","c"]`,
			Op:       ".[1]",
			Expected: `"b"`,
		},
		"range": {
			In:       `["a","b","c"]`,
			Op:       ".[1:2]",
			Expected: `["b","c"]`,
		},
		"nested index": {
			In:       `{"abc":"-","def":["a","b","c"]}`,
			Op:       ".def.[1]",
			Expected: `"b"`,
		},
		"nested range": {
			In:       `{"abc":"-","def":["a","b","c"]}`,
			Op:       ".def.[1:2]",
			Expected: `["b","c"]`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			op, err := jq.Parse(tc.Op)
			if err != nil {
				t.FailNow()
			}

			data, err := op.Apply([]byte(tc.In))
			if tc.HasError {
				if err == nil {
					t.FailNow()
				}
			} else {
				if string(data) != tc.Expected {
					t.FailNow()
				}
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}
