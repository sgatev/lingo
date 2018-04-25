package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestFuncCycloChecker(t *testing.T) {
	tests := []struct {
		description string
		max         int
		input       string
		expected    Report
	}{
		{
			description: "func, if statement, below limit",
			max:         3,
			input: `package test

			func foo(x int) int {
				if x > 5 {
					return 7
				}

				if x < 2 {
					return -7
				}

				return 21
			}
			`,
			expected: Report{},
		},
		{
			description: "method, if statement, below limit",
			max:         3,
			input: `package test

			type Foo struct{}

			func (f *Foo) foo(x int) int {
				if x > 5 {
					return 7
				}

				if x < 2 {
					return -7
				}

				return 21
			}
			`,
			expected: Report{},
		},
		{
			description: "func, if statement, above limit",
			max:         3,
			input: `package test

			func Bar(x int) int {
				if x > 5 {
					return 7
				}

				if x < 2 {
					return -7
				}

				if x > 43 {
					return 21
				}

				if x < -18 {
					return -21
				}

				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func Bar has cyclomatic complexity 5, max is 3",
					},
				},
			},
		},
		{
			description: "method, if statement, above limit",
			max:         3,
			input: `package test

			type Foo struct{}

			func (f *Foo) Bar(x int) int {
				if x > 5 {
					return 7
				}

				if x < 2 {
					return -7
				}

				if x > 43 {
					return 21
				}

				if x < -18 {
					return -21
				}

				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     40,
						Message: "func Bar has cyclomatic complexity 5, max is 3",
					},
				},
			},
		},
		{
			description: "method, nested statement, above limit",
			max:         3,
			input: `package test

			type Foo struct{}

			func (f *Foo) Bar(x int) int {
				if x > 5 {
					switch x {
					case 13:
						return 14
					case 21:
						return -5
					}
				}

				return 21
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     40,
						Message: "func Bar has cyclomatic complexity 4, max is 3",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewFuncCycloChecker(FuncCycloConfig{
				Max: test.max,
			}))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
