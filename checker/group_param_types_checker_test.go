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
			description: "func, no params",
			input: `
				package test

				func foo() {}
			`,
			expected: Report{},
		},
		{
			description: "func, single param",
			input: `
				package test

				func foo(a interface{}) {}
			`,
			expected: Report{},
		},
		{
			description: "func, ellipsis param",
			input: `
				package test

				func foo(a ...interface{}) {}
			`,
			expected: Report{},
		},
		{
			description: "func, different params, simple types",
			input: `
				package test

				func foo(a string, b int) {}
			`,
			expected: Report{},
		},
		{
			description: "func, same params, simple types, grouped",
			input: `
				package test

				func foo(a, b string) {}
			`,
			expected: Report{},
		},
		{
			description: "func, same params, simple types, not grouped",
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
			description: "method, same params, simple types, not grouped",
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
			description: "func, same params, pointer types, not grouped",
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
			description: "func, different params, array types",
			input: `
				package test

				func foo(a [3]int, b [5]int) {}
			`,
			expected: Report{},
		},
		{
			description: "func, same params, array types, not grouped",
			input: `
				package test

				func foo(a [3]int, b [3]int) {}
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
			description: "func, same params, slice types, not grouped",
			input: `
				package test

				func foo(a []int, b []int) {}
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
			description: "func, same params, map types, not grouped",
			input: `
				package test

				func foo(a map[string]int, b map[string]int) {}
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
			description: "func, same params, custom types, not grouped",
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
			description: "func, same params, func types, not grouped",
			input: `
				package test

				func foo(a func() int, b func() int) {}
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
			description: "func, different params, func types",
			input: `
				package test

				func foo(a func() int, b func() string) {}
			`,
			expected: Report{},
		},
		{
			description: "func, same params, chan types, not grouped",
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
			description: "func, different params, in-chan type",
			input: `
				package test

				func foo(a chan<- int, b chan int) {}
			`,
			expected: Report{},
		},
		{
			description: "func, different params, out-chan types",
			input: `
				package test

				func foo(a <-chan int, b chan int) {}
			`,
			expected: Report{},
		},
		{
			description: "func, same param types, separated by different one",
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
