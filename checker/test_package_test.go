package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestTestPackageChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "tests in non-test package",
			input: `
				package foo

				import "testing"

				func TestFoo(t *testing.t) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     6,
						Message: "package 'foo' should be named 'foo_test'",
					},
				},
			},
		},
		{
			description: "tests in test package",
			input: `
				package foo_test

				import "testing"

				func TestFoo(t *testing.t) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "non-tests in test package",
			input: `
				package foo_test

				func TestFoo(t *testing.t) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewTestPackageChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
