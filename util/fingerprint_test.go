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

import "testing"

type fingerprintTest struct {
	v interface{}
	e string
}

func TestFingerprint(t *testing.T) {

	tests := []fingerprintTest{
		fingerprintTest{
			v: struct {
				V0 string
				V1 string
			}{V0: "test", V1: "new"},
			e: "0e53a55537e3d923717949359b5e06edf02365f10633c8459778403c0a88dcad",
		},
		fingerprintTest{
			v: struct {
				V0 string
				V1 string
			}{V0: "diff", V1: "old"},
			e: "57f9eed2063446868f510626f89581e12d1eff0e57f17666b72251eaa90fe441",
		},
		fingerprintTest{
			v: struct {
				V0 string
				V2 string
			}{V0: "another", V2: "old"},
			e: "83d7fafa79795b64317d37faf2101bef19b601650c65ad16ea55ce52c50e5293",
		},
	}

	for _, test := range tests {
		r, err := Fingerprint(test.v)
		if err != nil {
			t.Fatalf("Fingerprint is broken %v", err)
		}

		if r != test.e {
			t.Errorf("Fingerprint is broken - received: %s - expected: %s", r, test.e)
		}
	}
}
