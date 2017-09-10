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

func BenchmarkObject(t *testing.B) {
	data := []byte(`{"hello":"world"}`)

	for i := 0; i < t.N; i++ {
		end, err := scanner.Object(data, 0)
		if err != nil {
			t.FailNow()
			return
		}

		if end == 0 {
			t.FailNow()
			return
		}
	}
}

func TestObject(t *testing.T) {
	testCases := map[string]struct {
		In     string
		Out    string
		HasErr bool
	}{
		"simple": {
			In:  `{"hello":"world"}`,
			Out: `{"hello":"world"}`,
		},
		"empty": {
			In:  `{}`,
			Out: `{}`,
		},
		"spaced": {
			In:  ` { "hello" : "world" } `,
			Out: ` { "hello" : "world" }`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			end, err := scanner.Object([]byte(tc.In), 0)
			if tc.HasErr {
				if err == nil {
					t.FailNow()
				}
			} else {
				data := tc.In[0:end]
				if string(data) != tc.Out {
					t.FailNow()
				}
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}
