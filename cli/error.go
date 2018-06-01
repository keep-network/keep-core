package cli

import (
	"fmt"
)

type (
	warningsType []string

	errWarnings struct {
		label    string
		warnings warningsType
	}
)

func (e *errWarnings) Error() string {
	// Populate label, warnings and error
	var warnings string
	for _, warning := range e.warnings {
		warnings += fmt.Sprintf("* %s\n", warning)
	}
	return fmt.Sprintf("%s:\n%s", e.label, warnings)
}

func (e *errWarnings) Print() {
	// Print error label and warnings
	fmt.Println(e.Error())
}
