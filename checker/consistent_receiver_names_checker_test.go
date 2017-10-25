package checker_test

import (
	"fmt"
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestConsistentReceiverNamesChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "inconsistent receiver names",
			input: `
				package foo

				func (f *Foo) Foo() {}
				func (b *Foo) Bar() {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("receivers in methods for type 'Foo' should have the same names"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(&ConsistentReceiverNamesChecker{})

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
