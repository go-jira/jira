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

package scanner

import (
	"fmt"
	"testing"
)

func TestSkipSpace(t *testing.T) {
	content := []byte(" \t\n\r!")
	end, err := skipSpace(content, 0)
	if err != nil {
		t.FailNow()
	}
	if end+1 != len(content) {
		t.FailNow()
	}
}

func TestExpect(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Expected string
		HasError bool
	}{
		"simple": {
			In:       "abc",
			Expected: "abc",
		},
		"extra": {
			In:       "abcdef",
			Expected: "abc",
		},
		"no match": {
			In:       "abc",
			Expected: "def",
			HasError: true,
		},
		"unexpected EOF": {
			In:       "ab",
			Expected: "abc",
			HasError: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			pos, err := expect([]byte(tc.In), 0, []byte(tc.Expected)...)
			if tc.HasError {
				if err == nil {
					t.FailNow()
				}

			} else {
				if err != nil {
					fmt.Println(err)
					t.FailNow()
				}
				if pos != len([]byte(tc.Expected)) {
					t.FailNow()
				}

			}
		})
	}
}
