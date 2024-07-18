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

package config

import "testing"

func TestConvertStrIntOrFatal(t *testing.T) {
	if v := convertStrIntOrFatal("10", "Test"); v != 10 {
		t.Errorf("convertStrIntOrFatal is broken - received: %d - expected: %d", 10, v)
	}

	if v := convertStrIntOrFatal("100", "Test"); v != 100 {
		t.Errorf("convertStrIntOrFatal is broken - received: %d - expected: %d", 100, v)
	}
}
