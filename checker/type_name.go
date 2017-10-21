package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

// TypeNameChecker checks the correctness of type names.
// Correct type names adhere to the following rules:
// * PascalCase for exported types.
// * camelCase for non-exported types.
type TypeNameChecker struct{}

// Check implements the Checker interface.
func (c *TypeNameChecker) Check(file *ast.File) error {
	for _, decl := range file.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if decl.Tok != token.TYPE {
			continue
		}
		spec, _ := decl.Specs[0].(*ast.TypeSpec)

		name := spec.Name.Name

		if !exportedNameRegexp.MatchString(name) &&
			!nonExportedNameRegexp.MatchString(name) {
			return fmt.Errorf("name '%s' is not valid", name)
		}
	}

	return nil
}

var (
	exportedNameRegexp    = regexp.MustCompile(`^([A-Z][a-z]*)+$`)
	nonExportedNameRegexp = regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)*$`)
)
