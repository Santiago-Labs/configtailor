package internal

import "fmt"

func MergeMaps(m1, m2 map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range m1 {
		merged[k] = v
	}

	for k, v := range m2 {
		merged[k] = v
	}

	return merged
}

func Flatten(m map[string]interface{}, prefix, connector string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		switch t := v.(type) {
		case map[string]interface{}:
			for innerKey, innerValue := range Flatten(t, k, "_") {
				result[innerKey] = innerValue
			}
		default:
			result[prefix+connector+k] = fmt.Sprintf("%s", v)
		}
	}
	return result
}
