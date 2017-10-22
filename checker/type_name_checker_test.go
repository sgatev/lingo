package checker_test

import (
	"fmt"
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestTypeName(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "correct exported",
			input: `
		package test

		type FooBar struct{}
	`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "correct non-exported",
			input: `
		package test

		type fooBar struct{}
	`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "incorrect",
			input: `
		package test

		type foo_bar struct{}
	`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("name 'foo_bar' is not valid"),
				},
			},
		},
		{
			description: "correct exported decl in func",
			input: `
		package test

		func foo() {
			type FooBar struct{}
		}
	`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "correct non-exported decl in func",
			input: `
		package test

		func foo() {
			type fooBar struct{}
		}
	`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "incorrect decl in func",
			input: `
		package test

		func foo() {
			type foo_bar struct{}
		}
	`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("name 'foo_bar' is not valid"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			fileChecker := NewFileChecker()
			checker := &TypeNameChecker{}
			checker.Register(fileChecker)

			file := ParseFileContent(test.input)
			var report Report
			fileChecker.Check(file, &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
