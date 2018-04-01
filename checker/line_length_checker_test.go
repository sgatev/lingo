package checker_test

import (
	"go/token"
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
						Pos:     19,
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

			fileSet := token.NewFileSet()
			file := ParseFileContentInSet(fileSet, test.input)
			var report Report
			checker.Check(file, test.input, &report)
			assert.Equal(t, test.expected, report)
			assert.Equal(t, 4, fileSet.Position(report.Errors[0].Pos).Line)
		})
	}
}
