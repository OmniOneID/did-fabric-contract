// Copyright 2024 Raonsecure

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
