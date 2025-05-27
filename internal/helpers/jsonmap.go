package helpers

import "gorm.io/datatypes"

func ToJSONMap(m map[string]int) datatypes.JSONMap {
	result := make(datatypes.JSONMap)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func FromJSONMap(m datatypes.JSONMap) map[string]int {
	result := make(map[string]int)
	for k, v := range m {
		if intVal, ok := v.(float64); ok {
			result[k] = int(intVal)
		}
	}
	return result
}
