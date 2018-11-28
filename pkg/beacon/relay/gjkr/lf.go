package gjkr

import (
	"fmt"
	"runtime"

	"github.com/pschlump/json" // modified from "encoding/json" to handle undefined types by ignoring them.
)

// LF Returns the File name and Line no as a string.
func LF(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("File: %s LineNo:%d", file, line)
	} else {
		return fmt.Sprintf("File: Unk LineNo:Unk")
	}
}

// ConvertToJSON return the JSON encoded version of the data with tab indentation.
func ConvertToJSON(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}
