package unique

//  Slice removes duplicate elements from a slice.
func Slice(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
