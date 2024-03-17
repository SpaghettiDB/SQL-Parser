package sqlParser 


// func Contains[T any](arr []any, element int) bool {
func ContainsAll(arr []string, elements []string) (bool, string) {
	for _, element := range elements {
		found := false
		for _, value := range arr {
			if value == element {
				found = true
				break
			}
		}
		if !found {
			return false, element
		}
	}
	return true, ""
}

