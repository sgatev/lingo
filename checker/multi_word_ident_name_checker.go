package checker

import (
	"fmt"
	"go/ast"
	"regexp"
)

func init() {
	must(Register("multi_word_ident_name", NewMultiWordIdentNameChecker))
}

// MultiWordIdentNameChecker checks the correctness of type names.
// Correct type names adhere to the following rules:
// * PascalCase for exported types.
// * camelCase for non-exported types.
type MultiWordIdentNameChecker struct{}

// NewMultiWordIdentNameChecker constructs a MultiWordIdentNameChecker.
func NewMultiWordIdentNameChecker(configData interface{}) NodeChecker {
	return &MultiWordIdentNameChecker{}
}

// Title implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Title() string {
	return "Multi-Word Identifiers"
}

// Description implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Description() string {
	return `An identifier consisting of multiple words must be in camelCase form.`
}

// Examples implements the NodeChecker interface.
func (c *MultiWordIdentNameChecker) Examples() []Example {
	return []Example{
		{
			Good: `type processTracker struct{}`,
			Bad:  `type process_tracker struct{}`,
		},
		{
			Good: `type ProcessTracker struct{}`,
			Bad:  `type Process_Tracker struct{}`,
		},
	}
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

	report.Errors = append(report.Errors, Error{
		Pos:     node.Pos(),
		Message: fmt.Sprintf("name '%s' is not valid", name),
	})
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
