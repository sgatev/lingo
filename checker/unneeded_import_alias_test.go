package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestUnneededImportAliasChecker(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    Report
	}{
		{
			description: "unique import, no alias",
			input: `
				package test

				import "math/rand"
			`,
			expected: Report{},
		},
		{
			description: "unique import, dot alias",
			input: `
				package test

				import . "math/rand"
			`,
			expected: Report{},
		},
		{
			description: "unique import, underscore alias",
			input: `
				package test

				import _ "math/rand"
			`,
			expected: Report{},
		},
		{
			description: "unique import, alias",
			input: `
				package test

				import something "math/rand"
			`,
			expected: Report{
				Errors: []Error{
					{
						Pos:     31,
						Message: "unneeded package alias: something",
					},
				},
			},
		},
		{
			description: "non-unique import, dot alias",
			input: `
				package test

				import (
					"math/rand"
					. "crypto/rand"
				)
			`,
			expected: Report{},
		},
		{
			description: "non-unique import, underscore alias",
			input: `
				package test

				import (
					"math/rand"
					_ "crypto/rand"
				)
			`,
			expected: Report{},
		},
		{
			description: "non-unique import, alias",
			input: `
				package test

				import (
					"math/rand"
					something "crypto/rand"
				)
			`,
			expected: Report{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			checker := NewFileChecker()
			checker.Register(NewUnneededImportAliasChecker(nil))

			file := ParseFileContent(test.input)
			var report Report
			checker.Check(file, "", &report)
			assert.Equal(t, test.expected, report)
		})
	}
}
