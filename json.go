// Copyright 2022 Stock Parfait

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"encoding/json"
)

// JSON parses js as a JSON string into the default encoding/json data
// structures: maps, strings, numbers, etc.  This is useful e.g. for custom JSON
// readers.
func JSON(js string) interface{} {
	var res interface{}
	if err := json.Unmarshal([]byte(js), &res); err != nil {
		panic(err)
	}
	return res
}
