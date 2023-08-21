package null

// All returns true if all given values are nil.
func All(values ...any) bool {
	for _, v := range values {
		if v != nil {
			return false
		}
	}
	return true
}
