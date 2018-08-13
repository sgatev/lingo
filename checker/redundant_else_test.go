package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestRedundantElseChecker(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    Report
	}{
		{
			description: "no else",
			input: `
				package test

				func do() int {
					if true {
						return 1
					}
				}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "else without terminating statement",
			input: `
				package test

				func do() {
					if true {
						foo()
					} else {
						bar()
					}
				}
			`,
			expected: Report{
				Errors: nil,
			},
		},
		{
			description: "return statement before else",
			input: `
				package test

				func do() int {
					if true {
						return 1
					} else {
						return 2
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     82,
						Message: "unexpected else after return statement",
					},
				},
			},
		},
		{
			description: "break statement before else",
			input: `
				package test

				func do() int {
					for {
						if true {
							break
						} else {
							return 2
						}
					}
					return 1
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     93,
						Message: "unexpected else after break statement",
					},
				},
			},
		},
		{
			description: "continue statement before else",
			input: `
				package test

				func do() int {
					for {
						if true {
							continue
						} else {
							return n
						}
					}
					return 1
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     96,
						Message: "unexpected else after continue statement",
					},
				},
			},
		},
		{
			description: "os.Exit() statement before else",
			input: `
				package test

				func do() int {
					if true {
						os.Exit(1)
					} else {
						return 2
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     84,
						Message: "unexpected else after os.Exit() statement",
					},
				},
			},
		},
		{
			description: "panic() statement before else",
			input: `
				package test

				func do() int {
					if true {
						panic("oops")
					} else {
						return 2
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     87,
						Message: "unexpected else after panic() statement",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewRedundantElseChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
