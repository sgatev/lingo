package checker

import (
	"fmt"
	"go/ast"
	"strings"
)

func init() {
	must(Register(NewLineLengthChecker))
}

// LineLengthChecker checks that code lines are withing specific length limits.
type LineLengthChecker struct {
	maxLineLength int
	tabWidth      int
}

// NewLineLengthChecker constructs a LineLengthChecker.
func NewLineLengthChecker() NodeChecker {
	return &LineLengthChecker{
		maxLineLength: 80,
		tabWidth:      4,
	}
}

// Slug implements the NodeChecker interface.
func (c *LineLengthChecker) Slug() string {
	return "line_length"
}

// Register implements the NodeChecker interface.
func (c *LineLengthChecker) Register(fc *FileChecker) {
	fc.On(&ast.File{}, c)
}

// Check implements the NodeChecker interface.
func (c *LineLengthChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	tabAsSpaces := strings.Repeat(" ", c.tabWidth)

	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		line = strings.Replace(line, "\t", tabAsSpaces, -1)

		if len(line) > c.maxLineLength {
			report.Errors = append(report.Errors,
				fmt.Errorf("line %d is too long", idx+1))
		}
	}
}
