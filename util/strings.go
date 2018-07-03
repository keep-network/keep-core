package util

import (
	"strings"
)

const (
	// MaxUint is the maximum unsigned integer value.
	MaxUint = ^uint(0)
	// MaxInt is the maximum integer value.
	MaxInt = int(MaxUint >> 1)
	// MinInt is the minimum integer value.
	MinInt = -MaxInt - 1
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
			frequencyMap[item]++
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

// Join concatenate the items in a list with the delimiter.
func Join(list []string, delimiter string) string {
	return strings.Join(list[:], delimiter)
}
