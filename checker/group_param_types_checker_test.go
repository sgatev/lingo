package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestGroupParamTypesChecker(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    Report
	}{
		{
			description: "no parameters",
			input: `
				package test

				func foo() {}
			`,
			expected: Report{},
		},
		{
			description: "single parameter",
			input: `
				package test

				func foo(a interface{}) {}
			`,
			expected: Report{},
		},
		{
			description: "ellipsis parameter",
			input: `
				package test

				func foo(a ...interface{}) {}
			`,
			expected: Report{},
		},
		{
			description: "different param types",
			input: `
				package test

				func foo(a string, b int) {}
			`,
			expected: Report{},
		},
		{
			description: "same param types, grouped",
			input: `
				package test

				func foo(a, b string) {}
			`,
			expected: Report{},
		},
		{
			description: "same simple param types, not grouped",
			input: `
				package test

				func foo(a string, b string) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: `params should be grouped by type`,
					},
				},
			},
		},
		{
			description: "method, same simple param types, not grouped",
			input: `
				package test

				func (f *Foo) Foo(a string, b string) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: `params should be grouped by type`,
					},
				},
			},
		},
		{
			description: "same pointer param types, not grouped",
			input: `
				package test

				func foo(a *string, b *string) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: `params should be grouped by type`,
					},
				},
			},
		},
		{
			description: "same type param types, not grouped",
			input: `
				package test

				func foo(a bar.Bar, b bar.Bar) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: `params should be grouped by type`,
					},
				},
			},
		},
		{
			description: "same type chan types, not grouped",
			input: `
				package test

				func foo(a chan int, b chan int) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: `params should be grouped by type`,
					},
				},
			},
		},
		{
			description: "same type in-chan types, grouped",
			input: `
				package test

				func foo(a chan<- int, b chan int) {}
			`,
			expected: Report{},
		},
		{
			description: "same type out-chan types, grouped",
			input: `
				package test

				func foo(a <-chan int, b chan int) {}
			`,
			expected: Report{},
		},
		{
			description: "same param types, separated by different one",
			input: `
				package test

				func foo(a string, b int, c string) {}
			`,
			expected: Report{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewGroupParamTypesChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
