package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestFuncResultsCountChecker(t *testing.T) {
	tests := []struct {
		description string
		max         int
		input       string
		expected    Report
	}{
		{
			description: "no results",
			max:         3,
			input: `package test

			func foo() {
			}
			`,
			expected: Report{},
		},
		{
			description: "less results",
			max:         3,
			input: `package test

			func foo() (int, error) {
				return 21, nil
			}
			`,
			expected: Report{},
		},
		{
			description: "equal results",
			max:         3,
			input: `package test

			func foo() (int, string, error) {
				return 21, "", nil
			}
			`,
			expected: Report{},
		},
		{
			description: "more results",
			max:         3,
			input: `package test

			func foo() (int, string, bool, error) {
				return 21, "", true, nil
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func has 4 results, max is 3",
					},
				},
			},
		},
		{
			description: "more named results",
			max:         3,
			input: `package test

			func foo() (a, b int, c, d bool) {
				return 1, 0, true, false
			}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     18,
						Message: "func has 4 results, max is 3",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewFuncResultsCountChecker(FuncResultsCountConfig{
				Max: test.max,
			}))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
