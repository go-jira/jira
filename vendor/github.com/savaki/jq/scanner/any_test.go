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

func BenchmarkAny(t *testing.B) {
	data := []byte(`"Hello, 世界 - 生日快乐"`)

	for i := 0; i < t.N; i++ {
		end, err := scanner.Any(data, 0)
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

func TestAny(t *testing.T) {
	testCases := map[string]struct {
		In  string
		Out string
	}{
		"string": {
			In:  `"hello"`,
			Out: `"hello"`,
		},
		"array": {
			In:  `["a","b","c"]`,
			Out: `["a","b","c"]`,
		},
		"object": {
			In:  `{"a":"b"}`,
			Out: `{"a":"b"}`,
		},
		"number": {
			In:  `1.234e+10`,
			Out: `1.234e+10`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			end, err := scanner.Any([]byte(tc.In), 0)
			if err != nil {
				t.FailNow()
			}
			data := tc.In[0:end]
			if string(data) != tc.Out {
				t.FailNow()
			}
		})
	}
}
