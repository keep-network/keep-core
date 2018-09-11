package jsonSyntaxErrorLib

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pschlump/godebug"
)

// HintType is a single hint on what might have cause the error
type HintType struct {
	Pattern string
	Note    string
	re      *regexp.Regexp
}

// HintList a set of possible, sometimes naive, hints on how to find/fix error
var HintList []HintType

// Debug if true will turn on extra debugging output
var Debug *bool

func init() {
	HintList = make([]HintType, 0, 25)
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key:value", Note: "Check for missing comma(,) or colon(:) immediately preceding this"})
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key$", Note: "Check for missing colon(:) betewen key and value"})
	HintList = append(HintList, HintType{Pattern: "unexpected end of JSON input", Note: "Check for missing end brace(}) or end array(])"})
	HintList = append(HintList, HintType{Pattern: "invalid character '\\\\'' looking for beginning of object key string", Note: "JSON only allows doublequotes(\")"})
	HintList = append(HintList, HintType{Pattern: "invalid character ',' looking for beginning of value", Note: "Check for two commas in a row, may be a missing value"})
	HintList = append(HintList, HintType{Pattern: "invalid character ']' looking for beginning of value", Note: "Check for an extra comma at end of array"})
	for ii, vv := range HintList {
		HintList[ii].re = regexp.MustCompile(vv.Pattern)
	}
	v := false
	Debug = &v
}

var hasTabs *regexp.Regexp

func init() {
	hasTabs = regexp.MustCompile("\t")
}

// GenerateSyntaxError converts from the offset error message into a human readable syntax error
func GenerateSyntaxError(js string, err error) (rv string) {

	max := func(a, b int) int {
		if a < b {
			return b
		}
		return a
	}
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// fmt.Printf("Type is [%T]\n", err)

	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		rv = fmt.Sprintf("%s\n", err)
		return
	}

	var sHint []string
	sErr := fmt.Sprintf("%s", err)
	for _, vv := range HintList {
		if vv.re.MatchString(sErr) {
			sHint = append(sHint, vv.Note)
		}
	}

	start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	start = min(max(start, 0), len(js)-1)
	end = min(max(end, 1), len(js))
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	if *Debug {
		fmt.Printf("AT: %s - start=%d end=%d len(js)=%d\n", godebug.LF(), start, end, len(js))
	}

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	pos = max(pos, 0)

	rv += fmt.Sprintf("Error in line %d: %s \n", line, err)
	rv += fmt.Sprintf("%s\n%s^\n", js[start:end], strings.Repeat(" ", pos))
	for _, vv := range sHint {
		// fmt.Printf("%s\n", vv)
		rv += vv + "\n"
	}
	return
}

// CheckForTabs returns true if data has tabs in it
func CheckForTabs(data []byte) bool {
	if hasTabs.MatchString(string(data)) {
		return true
	}
	return false
}

// TabListing shows tabs as '\\t' instead of a whitepace
func TabListing(data []byte) (rv string) {
	lineNo := 1
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		if hasTabs.MatchString(s) {
			s = strings.Replace(s, "\t", "\\t", -1)
			rv += fmt.Sprintf("%3d: %s\n", lineNo, s)
		}
		lineNo++
	}
	return
}
