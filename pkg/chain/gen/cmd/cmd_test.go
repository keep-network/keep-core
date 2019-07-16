package cmd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestComposableArgCheckerComposesOnSuccess(t *testing.T) {
	var seenFunctions []int

	firstComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 1)
		return nil
	})
	secondComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 2)
		return nil
	})
	thirdComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 3)
		return nil
	})

	tests := map[string]struct {
		composed composableArgChecker
		sequence []int
	}{
		"1-2-3": {
			composed: firstComposer.andThen(secondComposer).andThen(thirdComposer),
			sequence: []int{1, 2, 3},
		},
		"3-2-1": {
			composed: thirdComposer.andThen(secondComposer).andThen(firstComposer),
			sequence: []int{3, 2, 1},
		},
		"2-1-3": {
			composed: secondComposer.andThen(firstComposer).andThen(thirdComposer),
			sequence: []int{2, 1, 3},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			seenFunctions = []int{}
			result := test.composed(nil)

			if !reflect.DeepEqual(seenFunctions, test.sequence) {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.sequence,
					seenFunctions,
				)
			}

			if result != nil {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					nil,
					result,
				)
			}
		})
	}
}

func TestComposableArgCheckerShortCircuitsOnFailure(t *testing.T) {
	var seenFunctions []int

	firstComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 1)
		return nil
	})
	secondComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 2)
		return nil
	})
	thirdComposer := composableArgChecker(func(c *cli.Context) error {
		seenFunctions = append(seenFunctions, 3)
		return nil
	})
	errorComposer := composableArgChecker(func(c *cli.Context) error {
		return fmt.Errorf("ohai")
	})

	tests := map[string]struct {
		composed     composableArgChecker
		sequence     []int
		errorMessage string
	}{
		"1-2-error-3": {
			composed:     firstComposer.andThen(secondComposer).andThen(errorComposer).andThen(thirdComposer),
			sequence:     []int{1, 2},
			errorMessage: "ohai",
		},
		"3-error-2-1": {
			composed:     thirdComposer.andThen(errorComposer).andThen(secondComposer).andThen(firstComposer),
			sequence:     []int{3},
			errorMessage: "ohai",
		},
		"2-1-error-3": {
			composed:     secondComposer.andThen(firstComposer).andThen(errorComposer).andThen(thirdComposer),
			sequence:     []int{2, 1},
			errorMessage: "ohai",
		},
		"3-1-2-error": {
			composed:     thirdComposer.andThen(firstComposer).andThen(secondComposer).andThen(errorComposer),
			sequence:     []int{3, 1, 2},
			errorMessage: "ohai",
		},
		"3-1-2": {
			composed:     thirdComposer.andThen(firstComposer).andThen(secondComposer),
			sequence:     []int{3, 1, 2},
			errorMessage: "",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			seenFunctions = []int{}
			err := test.composed(nil)

			if !reflect.DeepEqual(seenFunctions, test.sequence) {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.sequence,
					seenFunctions,
				)
			}

			message := ""
			if err != nil {
				message = err.Error()
			}

			if message != test.errorMessage {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.errorMessage,
					message,
				)
			}
		})
	}
}
