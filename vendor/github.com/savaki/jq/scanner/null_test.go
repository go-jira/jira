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

func BenchmarkNull(t *testing.B) {
	data := []byte("null")
	for i := 0; i < t.N; i++ {
		pos, err := scanner.Null(data, 0)
		if err != nil {
			t.FailNow()
		}
		if pos != 4 {
			t.FailNow()
		}
	}
}
