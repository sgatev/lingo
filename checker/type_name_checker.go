package checker

import (
	"fmt"
	"go/ast"
	"regexp"
)

// TypeNameChecker checks the correctness of type names.
// Correct type names adhere to the following rules:
// * PascalCase for exported types.
// * camelCase for non-exported types.
type TypeNameChecker struct{}

// Register implements the NodeChecker interface.
func (c *TypeNameChecker) Register(fc *FileChecker) {
	fc.On(&ast.TypeSpec{}, c)
}

// Check implements the NodeChecker interface.
func (c *TypeNameChecker) Check(node ast.Node, report *Report) {
	spec := node.(*ast.TypeSpec)
	name := spec.Name.Name

	if !exportedNameRegexp.MatchString(name) &&
		!nonExportedNameRegexp.MatchString(name) {

		report.Errors = append(report.Errors,
			fmt.Errorf("name '%s' is not valid", name))
	}
}

var (
	exportedNameRegexp    = regexp.MustCompile(`^([A-Z][a-z]*)+$`)
	nonExportedNameRegexp = regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)*$`)
)
