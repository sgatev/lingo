package checker

import (
	"fmt"
	"go/ast"
	"regexp"
)

func init() {
	must(Register(NewMultiWordIdentNameChecker))
}

// MultiWordIdentNameChecker checks the correctness of type names.
// Correct type names adhere to the following rules:
// * PascalCase for exported types.
// * camelCase for non-exported types.
type MultiWordIdentNameChecker struct{}

// NewMultiWordIdentNameChecker constructs a MultiWordIdentNameChecker.
func NewMultiWordIdentNameChecker() NodeChecker {
	return &MultiWordIdentNameChecker{}
}

// Slug implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Slug() string {
	return "multi_word_ident_name"
}

// Register implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Register(fc *FileChecker) {
	fc.On(&ast.Ident{}, c)
}

// Check implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	name := node.(*ast.Ident).Name
	if isCorrectIdentName(name) {
		return
	}

	report.Errors = append(report.Errors,
		fmt.Errorf("name '%s' is not valid", name))
}

func isCorrectIdentName(name string) bool {
	return name == "_" ||
		exportedNameRegexp.MatchString(name) ||
		nonExportedNameRegexp.MatchString(name)
}

var (
	exportedNameRegexp = regexp.MustCompile(
		`^([A-Z0-9][a-z0-9]*)+$`)
	nonExportedNameRegexp = regexp.MustCompile(
		`^[a-z0-9]+([A-Z0-9][a-z0-9]*)*$`)
)
