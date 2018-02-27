package cli

import (
	"fmt"
	"testing"
)

func Test_CommandLineInterface(t *testing.T) {

	tests := []struct {
		args []string //
	}{
		{
			args: []string{"cli", "version"},
		},
	}

	for ii, test := range tests {
		var cli CLI
		fmt.Printf("Running test %d\n", ii)
		cli.Run(test.args...)
	}
}
