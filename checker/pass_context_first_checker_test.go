package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestPassContextFirstChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "context is not first parameter",
			input: `
				package foo

				func Foo(one int, ctx context.Context) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     23,
						Message: "func 'Foo' should be passed context as first parameter",
					},
				},
			},
		},
		{
			description: "context is first parameter",
			input: `
				package foo

				func Foo(ctx context.Context, one int) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "context is not first parameter of method",
			input: `
				package foo

				type Foo struct {}

				func (f *Foo) Bar(one int, ctx context.Context) {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     47,
						Message: "func 'Bar' should be passed context as first parameter",
					},
				},
			},
		},
		{
			description: "no parameters",
			input: `
				package foo

				func Foo() {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewPassContextFirstChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
