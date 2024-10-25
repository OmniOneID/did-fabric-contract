// Copyright 2024 OmniOne.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utility

import (
	"encoding/json"
	"sort"
)

func SortJson(jsonData []byte) []byte {
	var jsonMap map[string]interface{}
	json.Unmarshal(jsonData, &jsonMap)

	sortedJsonMap := SortJsonKeys(jsonMap)
	sortedJsonData, _ := json.Marshal(sortedJsonMap)
	return sortedJsonData
}

func SortJsonKeys(data map[string]interface{}) map[string]interface{} {
	sortedMap := make(map[string]interface{})
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		value := data[k]
		if nestedMap, ok := value.(map[string]interface{}); ok {
			sortedMap[k] = SortJsonKeys(nestedMap)
		} else {
			sortedMap[k] = value
		}
	}
	return sortedMap
}
