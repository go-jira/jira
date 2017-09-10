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

package scanner_test

import (
	"testing"

	"github.com/savaki/jq/scanner"
)

func BenchmarkFindRange(t *testing.B) {
	data := []byte(`["a","b","c","d","e"]`)

	for i := 0; i < t.N; i++ {
		out, err := scanner.FindRange(data, 0, 1, 2)
		if err != nil {
			t.FailNow()
			return
		}

		if string(out) != `["b","c"]` {
			t.FailNow()
			return
		}
	}
}

func TestFindRange(t *testing.T) {
	testCases := map[string]struct {
		In       string
		From     int
		To       int
		Expected string
		HasErr   bool
	}{
		"simple": {
			In:       `["a","b","c","d","e"]`,
			From:     1,
			To:       2,
			Expected: `["b","c"]`,
		},
		"single": {
			In:       `["a","b","c","d","e"]`,
			From:     1,
			To:       1,
			Expected: `["b"]`,
		},
		"mixed": {
			In:       `["a",{"hello":"world"},"c","d","e"]`,
			From:     1,
			To:       1,
			Expected: `[{"hello":"world"}]`,
		},
		"ordering": {
			In:     `["a",{"hello":"world"},"c","d","e"]`,
			From:   1,
			To:     0,
			HasErr: true,
		},
		"out of bounds": {
			In:     `["a",{"hello":"world"},"c","d","e"]`,
			From:   1,
			To:     20,
			HasErr: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := scanner.FindRange([]byte(tc.In), 0, tc.From, tc.To)
			if tc.HasErr {
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
