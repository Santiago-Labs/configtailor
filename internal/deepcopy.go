package internal

// Taken from: https://stackoverflow.com/a/62454413/6565736
type CopyableMap map[string]interface{}
type CopyableSlice []interface{}

// DeepCopy will create a deep copy of this map. The depth of this
// copy is all inclusive. Both maps and slices will be considered when
// making the copy.
func (m CopyableMap) DeepCopy() map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range m {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result[k] = CopyableMap(mapvalue).DeepCopy()
			continue
		}

		// Handle slices
		slicevalue, isSlice := v.([]interface{})
		if isSlice {
			result[k] = CopyableSlice(slicevalue).DeepCopy()
			continue
		}

		result[k] = v
	}

	return result
}

// DeepCopy will create a deep copy of this slice. The depth of this
// copy is all inclusive. Both maps and slices will be considered when
// making the copy.
func (s CopyableSlice) DeepCopy() []interface{} {
	result := []interface{}{}

	for _, v := range s {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result = append(result, CopyableMap(mapvalue).DeepCopy())
			continue
		}

		// Handle slices
		slicevalue, isSlice := v.([]interface{})
		if isSlice {
			result = append(result, CopyableSlice(slicevalue).DeepCopy())
			continue
		}

		result = append(result, v)
	}

	return result
}
