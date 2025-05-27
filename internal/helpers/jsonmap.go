package helpers

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
)

func ToJSONMap(m map[string]int) datatypes.JSONMap {
	result := make(datatypes.JSONMap)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func FromJSONMap(m datatypes.JSONMap) (map[string]int, error) {
	result := make(map[string]int)
	for k, v := range m {
		switch val := v.(type) {
		case json.Number:
			intVal, err := val.Int64()
			if err == nil {
				result[k] = int(intVal)
			} else {
				return nil, fmt.Errorf("Invalid json.Number for key %s: %v\n", k, err)
			}
		case float64:
			result[k] = int(val)
		case string:
			var num json.Number = json.Number(val)
			if intVal, err := num.Int64(); err == nil {
				result[k] = int(intVal)
			}
		default:
			return nil, fmt.Errorf("Unknown type for key %s: %T\n", k, v)
		}
	}
	return result, nil
}
