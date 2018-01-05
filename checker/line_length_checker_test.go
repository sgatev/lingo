package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestLineLengthChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "long line",
			input: `
				package foo

				func TestFooBarFunctionVeryLong(a int, b int, c int, d int) (error, float64) {}`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     11,
						Message: "line is too long",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewLineLengthChecker(&LineLengthConfig{
				MaxLength: 80,
				TabWidth:  4,
			}))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, test.input, &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
