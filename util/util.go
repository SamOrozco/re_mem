package util

import (
	"strings"
	"unicode"
)

const delim = '_'
// This method is meant to cleanse the attribute Name
// and make it a viable "Key"
func CleanseName(val string) string {
	// trim spaces
	val = strings.TrimSpace(val)
	// the cleanse algorithm will iterate each char in the string
	// if it runs into a character that is not a number or letter
	// it will continue until it finds a letter and then replace the non letter or number void with
	// an underscore
	result := make([]rune, 0)
	lastInvalidChar := false
	for _, r := range val {
		r = unicode.ToLower(r)
		isValidChar := unicode.IsLetter(r) || unicode.IsNumber(r)
		if isValidChar {
			if lastInvalidChar {
				result = append(result, delim)
				result = append(result, r)
				lastInvalidChar = false
			} else {
				result = append(result, r)
			}
		} else {
			lastInvalidChar = true
		}
	}
	return string(result)
}
