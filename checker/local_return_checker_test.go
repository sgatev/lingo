package checker_test

import (
	"fmt"
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestLocalReturnChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "local func, local return",
			input: `
				package test

				func foo1() bar {}

				func (f *Foo) foo2() bar {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "exported func, exported return",
			input: `
				package test

				func Foo1() Bar {}

				func (f *Foo) Foo2() Bar {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "local func, exported return",
			input: `
				package test

				func foo1() Bar {}

				func (f *Foo) foo2() Bar {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "local func, mixed returns",
			input: `
				package test

				func foo1() (bar1, Bar2, bar3) {}

				func (f *Foo) foo2() (bar1, Bar2, bar3) {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "exported func, local return",
			input: `
				package test

				func Foo1() bar {}

				func (f *Foo) Foo2() bar {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar'"),
				},
			},
		},
		{
			description: "exported func, mixed returns",
			input: `
				package test

				func Foo1() (bar1, Bar2, bar3) {}

				func (f *Foo) Foo2() (bar1, Bar2, bar3) {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar1'"),
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar3'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar1'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar3'"),
				},
			},
		},
		{
			description: "exported func, internal return",
			input: `
				package test

				func Foo1() string {}

				func (f *Foo) Foo2() string {}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "exported func, chan of local return",
			input: `
				package test

				func Foo1() chan bar {}

				func (f *Foo) Foo2() chan bar {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar'"),
				},
			},
		},
		{
			description: "exported func, slice of local return",
			input: `
				package test

				func Foo1() []bar {}

				func (f *Foo) Foo2() []bar {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar'"),
				},
			},
		},
		{
			description: "exported func, array of local return",
			input: `
				package test

				func Foo1() [21]bar {}

				func (f *Foo) Foo2() [21]bar {}
			`,
			expected: Report{
				Errors: []error{
					fmt.Errorf("exported func 'Foo1' cannot return value of local type 'bar'"),
					fmt.Errorf("exported func 'Foo2' cannot return value of local type 'bar'"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			fileChecker := NewFileChecker()
			checker := &LocalReturnChecker{}
			checker.Register(fileChecker)

			file := ParseFileContent(test.input)
			var report Report
			fileChecker.Check(file, &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
