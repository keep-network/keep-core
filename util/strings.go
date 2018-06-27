package util

import (
	"fmt"
	"regexp"
)

// CompileRegex compiles the regex pattern for later use
func CompileRegex(pattern string) *regexp.Regexp {
	r, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("Error compiling regex: [%s]", pattern))
	}

	return r
}

// MatchFound returns true if the search term matches the pattern
func MatchFound(regexPattern *regexp.Regexp, term string) bool {
	return regexPattern.FindString(term) != ""
}
