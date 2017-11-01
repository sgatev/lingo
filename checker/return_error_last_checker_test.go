package checker_test

import (
	"fmt"
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestReturnErrorLastChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "error not last return",
			input: `
				package foo

				func Foo() (int, error, bool) {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("func 'Foo' should return error as the last value"),
				},
			},
		},
		{
			description: "error last return",
			input: `
				package foo

				func Foo() (int, bool, error) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "no return",
			input: `
				package foo

				func Foo() (int, bool, error) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "return not ident",
			input: `
				package foo

				func Foo() (struct{},error) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewReturnErrorLastChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
