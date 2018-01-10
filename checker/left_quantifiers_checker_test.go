package checker_test

import (
	"fmt"
	"testing"

	. "github.com/s2gatev/lingo/checker"
	"github.com/stretchr/testify/assert"
)

func TestLeftQuantifiers(t *testing.T) {
	type test struct {
		description string
		expression  string
		expected    Report
	}

	input := `
		package main

		import "time"

		func main() {
			%s
		}
	`

	tests := []test{
		{
			description: "not binary expression",
			expression:  `_ = time.Second(1000)`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "basic literal expression",
			expression:  `_ = 5 * 5`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "basic literal is not on the left",
			expression:  `_ = time.Second * 5`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     58,
						Message: "the left operand should be a basic literal",
					},
				},
			},
		},
		{
			description: "left quantifiers",
			expression:  `_ = 5 * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "no basic literals",
			expression:  `_ = time.Duration(1) * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "multiple left quantifiers",
			expression:  `_ = 2 + 5 * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "multiple non-basic expressions on the right",
			expression:  `_ = 2 * time.Duration(2) * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "multiple binary operations, basic literal at the end",
			expression:  `_ = 2 * time.Second * 3`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     58,
						Message: "the left operand should be a basic literal",
					},
				},
			},
		},
		{
			description: "binary expression with parentheses",
			expression:  `_ = (60 + 10) * 60 * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "mixed expressions",
			expression:  `_ = time.Duration(1) * 2 * 5 * time.Second`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     58,
						Message: "the left operand should be a basic literal",
					},
				},
			},
		},
		{
			description: "non-commutative operators",
			expression:  `_ = time.Duration(1) - 1`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "mixed operators",
			expression:  `_ =  time.Duration(1) * (i - 1)`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "binary operators",
			expression:  `_ = i & 1`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     58,
						Message: "the left operand should be a basic literal",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewLeftQuantifiersChecker(nil))

			file := ParseFileContent(fmt.Sprintf(input, test.expression))
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
