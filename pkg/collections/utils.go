package collections

// Delete from index in array preserving the order
func DeletePreserving(a []interface{}, index int) []interface{} {
	if index > len(a) || index < 0 {
		return nil
	}
	a, a[len(a)-1] = append(a[:index], a[index+1:]...), nil
	return a
}

// Add at index in array
func AddAt(a []interface{}, index int, value interface{}) []interface{} {
	if index > len(a) || index < 0 {
		return nil
	}
	tmp := append(a[index:], value)
	tmp = append(tmp, a[:index])
	return tmp
}

// Check if map with string keys contains all keys from the given array
func ArrayOfStringContains(sourceValues []string, values []string) bool {
	sourceMap := make(map[interface{}]struct{})
	for _, sourceValue := range sourceValues {
		sourceMap[sourceValue] = struct{}{}
	}

	for _, value := range values {
		if _, ok := sourceMap[value]; !ok {
			return false
		}
	}

	return true
}
