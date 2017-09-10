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

func BenchmarkChain(t *testing.B) {
	op := jq.Chain(jq.Dot("a"), jq.Dot("b"))
	data := []byte(`{"a":{"b":"value"}}`)

	for i := 0; i < t.N; i++ {
		_, err := op.Apply(data)
		if err != nil {
			t.FailNow()
			return
		}
	}
}

func TestChain(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Op       jq.Op
		Expected string
		HasError bool
	}{
		"simple": {
			In:       `{"hello":"world"}`,
			Op:       jq.Chain(jq.Dot("hello")),
			Expected: `"world"`,
		},
		"nested": {
			In:       `{"a":{"b":"world"}}`,
			Op:       jq.Chain(jq.Dot("a"), jq.Dot("b")),
			Expected: `"world"`,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := tc.Op.Apply([]byte(tc.In))
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
