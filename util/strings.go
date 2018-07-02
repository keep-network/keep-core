package util

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

// maxMapVal returns the maximum value found in the map or MinInt for an empty map.
func maxMapVal(freqMap map[string]int) int {
	maxVal := MinInt
	for _, val := range freqMap {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// frequencyMap returns a map with the count of each string in slice.
func frequencyMap(list []string) map[string]int {
	frequencyMap := make(map[string]int)
	for _, item := range list {
		_, exist := frequencyMap[item]
		if exist {
			frequencyMap[item] += 1
		} else {
			frequencyMap[item] = 1
		}
	}
	return frequencyMap
}

// DuplicatesExist returns true if one or more duplicates exist in the list.
func DuplicatesExist(list []string) bool {
	return maxMapVal(frequencyMap(list)) > 1
}

// Duplicates returns the list of duplicated items in list.
func Duplicates(list []string) []string {
	return duplicatesInMap(frequencyMap(list))
}

// duplicatesInMap returns the list of duplicated items in list.
func duplicatesInMap(freqMap map[string]int) []string {
	dups := []string{}
	for key, val := range freqMap {
		if val > 1 {
			dups = append(dups, key)
		}
	}
	return dups
}

// CompileRegex compiles the regex pattern for later use.
func CompileRegex(pattern string) *regexp.Regexp {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("Error compiling regex: [%s]", pattern))
	}
	return r
}

// MatchFound returns true if the search term matches the pattern.
func MatchFound(regexPattern *regexp.Regexp, term string) bool {
	return regexPattern.FindString(term) != ""
}

// Join concatenate the items in a list with the delimiter.
func Join(list []string, delimiter string) string {
	return strings.Join(list[:], delimiter)
}
