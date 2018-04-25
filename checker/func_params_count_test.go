package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestFuncParamsCountChecker(t *testing.T) {
	tests := []struct {
		description string
		max         int
		input       string
		expected    Report
	}{
		{
			description: "no params",
			max:         3,
			input: `package test

			func foo() int {
				return 21
			}
			`,
			expected: Report{},
		},
		{
			description: "less params",
			max:         3,
			input: `package test

			func foo(a int, b string) int {
				return 21
			}
			`,
			expected: Report{},
		},
		{
			description: "equal params",
			max:         3,
			input: `package test

			func foo(a int, b string, c ...float64) int {
				return 21
			}
			`,
			expected: Report{},
		},
		{
			description: "more params",
			max:         3,
			input: `package test

			func foo(a int, b, c string, d float64) int {
				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func has 4 params, max is 3",
					},
				},
			},
		},
		{
			description: "more unnamed params",
			max:         3,
			input: `package test

			func foo(int, string, string, float64) int {
				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func has 4 params, max is 3",
					},
				},
			},
		},
		{
			description: "more params with skipped names",
			max:         3,
			input: `package test

			func foo(_, _ int, a string, b float64) int {
				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func has 4 params, max is 3",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewFuncParamsCountChecker(FuncParamsCountConfig{
				Max: test.max,
			}))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
