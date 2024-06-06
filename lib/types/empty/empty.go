package empty

// Count returns how many empty values were provided.
func Count(values ...string) int {
	emptyValues := 0
	for _, v := range values {
		if len(v) == 0 {
			emptyValues++
		}
	}
	return emptyValues
}
