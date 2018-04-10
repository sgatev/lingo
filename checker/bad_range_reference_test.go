package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestBadRangeReferenceChecker(t *testing.T) {
	type test struct {
		description string
		input       string
		expected    Report
	}

	tests := []test{
		{
			description: "key value",
			input: `
						package test

						func foo(foos []string) {
							for foo := range foos {
								bar(foo)
							}
						}
					`,
			expected: Report{},
		},
		{
			description: "bad key reference",
			input: `
						package test

						func foo(foos []string) {
							for foo := range foos {
								bar(&foo)
							}
						}
					`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     97,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "bad value reference",
			input: `
						package test

						func foo(foos []string) {
							for _, foo := range foos {
								bar(&foo)
							}
						}
					`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     100,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "bad value reference, method call",
			input: `
						package test

						func foo(foos []string) {
							for _, foo := range foos {
								b.Bar(&foo)
							}
						}
					`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     102,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "overridden value reference",
			input: `
						package test

						func foo(foos []string) {
							for _, foo := range foos {
								foo := foo
								bar(&foo)
							}
						}
					`,
			expected: Report{},
		},
		{
			description: "bad value overridden after reference",
			input: `
						package test

						func foo(foos []string) {
							for _, foo := range foos {
								bar(&foo)
								foo := foo
								qux(&foo)
							}
						}
					`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     100,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "value reference via func",
			input: `
				package test

				func foo(foos []string) {
					for _, foo := range foos {
						func(foo string) {
							bar(&foo)
						}(foo)
					}
				}
			`,
			expected: Report{},
		},
		{
			description: "value reference via go func",
			input: `
				package test

				func foo(foos []string) {
					for _, foo := range foos {
						go func(foo string) {
							bar(&foo)
						}(foo)
					}
				}
			`,
			expected: Report{},
		},
		{
			description: "nested bad value reference via func",
			input: `
				package test

				func foo(foos []string) {
					for _, foo := range foos {
						func(foo string) {
							for _, foo := range foos {
								bar(&foo)
							}
						}(foo)
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     153,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "nested bad value reference via go func",
			input: `
				package test

				func foo(foos []string) {
					for _, foo := range foos {
						go func(foo string) {
							for _, foo := range foos {
								bar(&foo)
							}
						}(foo)
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     156,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "bad value reference, assignment",
			input: `
				package test

				func foo(foos []string) {
					for _, foo := range foos {
						qux := &foo
						bar(qux)
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     95,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
		{
			description: "bad value reference, append",
			input: `
				package test

				var res []*string

				func foo(foos []string) {
					for _, foo := range foos {
						res = append(res, &foo)
					}
				}
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     129,
						Message: "bad reference of range var: foo",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewBadRangeReferenceChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
