package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("line_length", NewLineLengthChecker))
}

// LineLengthConfig describes the configuration of a LineLengthChecker.
type LineLengthConfig struct {

	// MaxLength is the maximum number of characters permitted on a single line.
	MaxLength int `mapdecode:"max_length"`

	// TabWidth is the number of characters equivalent to a single tab.
	TabWidth int `mapdecode:"tab_width"`
}

// LineLengthChecker checks that code lines are within specific length limits.
type LineLengthChecker struct {
	maxLength int
	tabWidth  int
}

// NewLineLengthChecker constructs a LineLengthChecker.
func NewLineLengthChecker(configData interface{}) NodeChecker {
	var config LineLengthConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &LineLengthChecker{
		maxLength: config.MaxLength,
		tabWidth:  config.TabWidth,
	}
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

	pos := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.Replace(line, "\t", tabAsSpaces, -1)

		if len(line) > c.maxLength {
			report.Errors = append(report.Errors, Error{
				Pos:     token.Pos(pos),
				Message: fmt.Sprintf("line is too long"),
			})
		}

		pos += len(line)
	}
}
