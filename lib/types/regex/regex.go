package regex

import (
	"regexp"
)

// MatchAny matches a value against one or more regex patterns.
func MatchAny(value string, patterns ...string) bool {
	for _, pattern := range patterns {
		match, _ := regexp.MatchString(pattern, value)
		if match {
			return true
		}
	}
	return false
}
