package checker_test

import (
	"go/token"
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestLineLengthChecker(t *testing.T) {
	input := `package foo

	func TestFooBarFunctionVeryLong(a int, b int, c int, d int) (error, float64) {}`

	checker := NewFileChecker()
	checker.Register(NewLineLengthChecker(&LineLengthConfig{
		MaxLength: 80,
		TabWidth:  4,
	}))

	fileSet := token.NewFileSet()
	file := ParseFileContentInSet(fileSet, input)

	var report Report
	checker.Check(file, input, &report)

	assert.Equal(t,
		Report{
			Errors: []Error{
				{
					Pos:     14,
					Message: "line is too long",
				},
			},
		},
		report)
	assert.Equal(t, 3, fileSet.Position(report.Errors[0].Pos).Line)
}
