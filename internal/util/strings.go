package util

import "strings"

// NormalizeName normalizes a type or attribute name for lookup & matching
func NormalizeName(input string) string {
	return strings.ToLower(input)
}

// Longest returns the length of the longest string in the input slice
func Longest(input []string) int {
	longest := 0

	for _, s := range input {
		if l := len(s); l > longest {
			longest = l
		}
	}

	return longest
}
