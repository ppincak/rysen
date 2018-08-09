package collections

func DeletePreserving(a []interface{}, index int) []interface{} {
	if index > len(a) || index < 0 {
		return nil
	}
	a, a[len(a)-1] = append(a[:index], a[index+1:]...), nil
	return a
}

func AddAt(a []interface{}, index int, value interface{}) []interface{} {
	if index > len(a) || index < 0 {
		return nil
	}
	tmp := append(a[index:], value)
	tmp = append(tmp, a[:index])
	return tmp
}