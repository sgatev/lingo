package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestExportedIdentDocChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "type",
			input: `
				package test

				// FooBar1 is documented.
				type FooBar1 struct{}

				type fooBar2 struct{}

				type FooBar3 struct{}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     113,
						Message: "exported identifier 'FooBar3' is not documented",
					},
				},
			},
		},
		{
			description: "const",
			input: `
				package test

				// TheAnswer1 is documented.
				const TheAnswer1 = 42

				const theAnswer2 = 42

				const TheAnswer3 = 42
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     117,
						Message: "exported identifier 'TheAnswer3' is not documented",
					},
				},
			},
		},
		{
			description: "var",
			input: `
				package test

				// TheAnswer1 is documented.
				var TheAnswer1 = 42

				var theAnswer2 = 42

				var TheAnswer3 = 42
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     111,
						Message: "exported identifier 'TheAnswer3' is not documented",
					},
				},
			},
		},
		{
			description: "func",
			input: `
				package test

				// Foo1 is documented.
				func Foo1() {}

				func foo2() {}

				func Foo3() {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     91,
						Message: "exported identifier 'Foo3' is not documented",
					},
				},
			},
		},
		{
			description: "struct field",
			input: `
				package test

				// Foo is documented.
				type Foo struct {

					// FooBar1 is documented.
					FooBar1	int
					fooBar2	int
					FooBar3	int

					FooBar4, fooBar5, FooBar6 int

					// Fields are documented.
					FooBar7, fooBar8, FooBar9 int
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     139,
						Message: "exported identifier 'FooBar3' is not documented",
					},
					{
						Pos:     157,
						Message: "exported identifier 'FooBar4' is not documented",
					},
					{
						Pos:     157,
						Message: "exported identifier 'FooBar6' is not documented",
					},
				},
			},
		},
		{
			description: "interface method",
			input: `
				package test

				// Foo is documented.
				type Foo interface {

					// FooBar1 is documented.
					FooBar1()	int
					fooBar2()	int
					FooBar3()	int
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     146,
						Message: "exported identifier 'FooBar3' is not documented",
					},
				},
			},
		},
		{
			description: "grouped identifiers",
			input: `
				package test

				var (

					// Foo is documented.
					Foo int

					// Bar is documented.
					Bar string
				)
			`,
			expected: Report{
				Errors: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewExportedIdentDocChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}

func TestExportedIdentDocCheckerPrefix(t *testing.T) {
	config := &ExportedIdentDocCheckerConfig{
		HasIdentPrefix: true,
	}

	tests := []struct {
		description string
		input       string
		expected    Report
	}{
		{
			description: "type",
			input: `
				package test

				// Some documentation.
				type FooBar1 struct{}

				// FooBar2 is documented.
				type FooBar2 struct{}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: "expected the comment to start with 'FooBar1'",
					},
				},
			},
		},
		{
			description: "var",
			input: `
				package test

				// Some documentation.
				var FooBar1 int

				// FooBar2 is documented.
				var FooBar2 string

				var FooBar3 string
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: "expected the comment to start with 'FooBar1'",
					},
					{
						Pos:     130,
						Message: "exported identifier 'FooBar3' is not documented",
					},
				},
			},
		},
		{
			description: "groped identifiers",
			input: `
				package test

				// Some documentation.
				var (
					// Some other documentation.
					FooBar1 int

					// FooBar2 is documented.
					FooBar2 string
				)
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     62,
						Message: "expected the comment to start with 'FooBar1'",
					},
				},
			},
		},
		{
			description: "func",
			input: `
				package test

				// Some documentation.
				func Foo1() {}

				// Foo2 is documented.
				func Foo2() {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: "expected the comment to start with 'Foo1'",
					},
				},
			},
		},
		{
			description: "fields",
			input: `
				package test

				// Some documentation.
				type Foo struct {

					// Some documentation.
					FooBar1	int

					// FooBar2 is documented.
					FooBar2	int

					// Fields are documented.
					FooBar7, fooBar8, FooBar9 int
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     24,
						Message: "expected the comment to start with 'Foo'",
					},
					{
						Pos:     75,
						Message: "expected the comment to start with 'FooBar1'",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewExportedIdentDocChecker(config))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
