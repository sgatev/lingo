package checker

import (
	"bufio"
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

// Title implements the NodeChecker interface.
func (c *LineLengthChecker) Title() string {
	return "Line Length"
}

// Description implements the NodeChecker interface.
func (c *LineLengthChecker) Description() string {
	return fmt.Sprintf(`The maximum line of a length is %d symbols.`, c.maxLength)
}

// Examples implements the NodeChecker interface.
func (c *LineLengthChecker) Examples() []Example {
	return nil
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

	pos := int(node.Pos())
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		length := len(line) + 1

		line = strings.Replace(line, "\t", tabAsSpaces, -1)
		if len(line) > c.maxLength {
			report.Errors = append(report.Errors, Error{
				Pos:     token.Pos(pos),
				Message: fmt.Sprintf("line is too long"),
			})
		}

		pos += length
	}
}
