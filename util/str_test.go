/**
 * Copyright 2024-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"testing"
	"time"
)

type strIntDurationTest struct {
	v  string
	dt time.Duration
	e  time.Duration
}

func TestConvertStrIntToDuration(t *testing.T) {

	tests := []strIntDurationTest{
		strIntDurationTest{
			v:  "10",
			dt: time.Second,
			e:  10 * time.Second,
		},
		strIntDurationTest{
			v:  "20",
			dt: time.Minute,
			e:  20 * time.Minute,
		},
	}

	for _, test := range tests {
		r, err := ConvertStrIntToDuration(test.v, test.dt)
		if err != nil {
			t.Fatalf("ConvertStrIntToDuration is broken %v", err)
		}

		if r != test.e {
			t.Errorf("ConvertStrIntToDuration is broken - received: %v - expected: %v", r, test.e)
		}
	}
}
