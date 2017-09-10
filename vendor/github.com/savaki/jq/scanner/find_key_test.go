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

func BenchmarkFindKey(t *testing.B) {
	data := []byte(`{"hello":"world"}`)

	for i := 0; i < t.N; i++ {
		out, err := scanner.FindKey(data, 0, []byte("hello"))
		if err != nil {
			t.FailNow()
			return
		}

		if string(out) != `"world"` {
			t.FailNow()
			return
		}
	}
}

func TestFindKey(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Key      string
		Expected string
		HasErr   bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Key:      "hello",
			Expected: `"world"`,
		},
		"spaced": {
			In:       ` { "hello" : "world" } `,
			Key:      "hello",
			Expected: `"world"`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := scanner.FindKey([]byte(tc.In), 0, []byte(tc.Key))
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
