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
			description: "basic literal is not on the left",
			expression:  `_ = time.Second * 5`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("the left operand should be a basic literal"),
				},
			},
		},
		{
			description: "multiple binary operations, basic literal at the end",
			expression:  `_ = 2 + time.Second * 3`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("the left operand should be a basic literal"),
				},
			},
		},
		{
			description: "mixed expressions",
			expression:  `_ = time.Duration(1) * 2 * 5 * time.Second`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("the left operand should be a basic literal"),
				},
			},
		},
		{
			description: "multiple binary operations, not a basic literal at the end",
			expression:  `_ = 60 - 10 + 60 * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "binary expression with parentheses",
			expression:  `_ = (60 - 10) * 60 * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "more than one non-basic expression on the right",
			expression:  `_ = 2 + 22 * time.Duration(1) * time.Second`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "no basic literal in binary expression",
			expression:  `_ = time.Duration(1) * time.Second`,
			expected: Report{
				Errors: nil,
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
