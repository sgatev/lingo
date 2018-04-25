package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestMultiWordIdentNameChecker(t *testing.T) {
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

				type FooBar1 struct{}
				type fooBar2 struct{}
				type foo_bar3 struct{}

				func foo() {
					type FooBar4 struct{}
					type fooBar5 struct{}
					type foo_bar6 struct{}

					_ = func() {
						type FooBar7 struct{}
						type fooBar8 struct{}
						type foo_bar9 struct{}
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     81,
						Message: "name 'foo_bar3' is not valid",
					},
					{
						Pos:     181,
						Message: "name 'foo_bar6' is not valid",
					},
					{
						Pos:     285,
						Message: "name 'foo_bar9' is not valid",
					},
				},
			},
		},
		{
			description: "const",
			input: `
				package test

				const TheAnswer1 = 42
				const theAnswer2 = 42
				const the_answer3 = 42

				func foo() {
					const TheAnswer4 = 42
					const theAnswer5 = 42
					const the_answer6 = 42

					_ = func() {
						const TheAnswer7 = 42
						const theAnswer8 = 42
						const the_answer9 = 42
					}

					const FooBar, fooBar, foo_bar = 1, 2, 3
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     82,
						Message: "name 'the_answer3' is not valid",
					},
					{
						Pos:     182,
						Message: "name 'the_answer6' is not valid",
					},
					{
						Pos:     286,
						Message: "name 'the_answer9' is not valid",
					},
					{
						Pos:     338,
						Message: "name 'foo_bar' is not valid",
					},
				},
			},
		},
		{
			description: "var",
			input: `
				package test

				var TheAnswer1 = 42
				var theAnswer2 = 42
				var the_answer3 = 42

				func foo() {
					var TheAnswer4 = 42
					var theAnswer5 = 42
					var the_answer6 = 42

					_ = func() {
						var TheAnswer7 = 42
						var theAnswer8 = 42
						var the_answer9 = 42

						FooBar4 := Foo()
						fooBar5 := Foo()
						foo_bar6 := Foo()
					}

					var FooBar1, fooBar2, foo_bar3 = Foo()
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     76,
						Message: "name 'the_answer3' is not valid",
					},
					{
						Pos:     170,
						Message: "name 'the_answer6' is not valid",
					},
					{
						Pos:     268,
						Message: "name 'the_answer9' is not valid",
					},
					{
						Pos:     338,
						Message: "name 'foo_bar6' is not valid",
					},
					{
						Pos:     391,
						Message: "name 'foo_bar3' is not valid",
					},
				},
			},
		},
		{
			description: "func",
			input: `
				package test

				func FooBar1() {}
				func fooBar2() {}
				func foo_bar3() {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     73,
						Message: "name 'foo_bar3' is not valid",
					},
				},
			},
		},
		{
			description: "struct method",
			input: `
				package test

				type Foo struct{}

				func (f *Foo) FooBar1() {}
				func (f *Foo) fooBar2() {}
				func (f *Foo) foo_bar3() {}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     123,
						Message: "name 'foo_bar3' is not valid",
					},
				},
			},
		},
		{
			description: "struct field",
			input: `
				package test

				type Foo struct {
					FooBar1		int
					fooBar2		int
					foo_bar3	int

					FooBar4, fooBar5, foo_bar6 int
				}

				func foo() {
					type Foo struct {
						FooBar7		int
						fooBar8		int
						foo_bar9	int

						FooBar10, fooBar11, foo_bar12 int
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     83,
						Message: "name 'foo_bar3' is not valid",
					},
					{
						Pos:     120,
						Message: "name 'foo_bar6' is not valid",
					},
					{
						Pos:     224,
						Message: "name 'foo_bar9' is not valid",
					},
					{
						Pos:     264,
						Message: "name 'foo_bar12' is not valid",
					},
				},
			},
		},
		{
			description: "interface method",
			input: `
				package test

				type Foo interface {
					FooBar1()	int
					fooBar2()	int
					foo_bar3()	int
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     88,
						Message: "name 'foo_bar3' is not valid",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewMultiWordIdentNameChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
