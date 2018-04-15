package checker

import (
	"fmt"
	"go/ast"
	"path"
	"strings"
)

func init() {
	must(Register("unneeded_import_alias", NewUnneededImportAliasChecker))
}

// UnneededImportAliasChecker checks that import aliases are used only when
// necessary.
type UnneededImportAliasChecker struct{}

// NewUnneededImportAliasChecker constructs a UnneededImportAliasChecker.
func NewUnneededImportAliasChecker(configData interface{}) NodeChecker {
	return &UnneededImportAliasChecker{}
}

// Title implements the NodeChecker interface.
func (c *UnneededImportAliasChecker) Title() string {
	return "Unneeded Import Alias"
}

// Description implements the NodeChecker interface.
func (c *UnneededImportAliasChecker) Description() string {
	return `Import aliases should be used only when necessary.`
}

// Examples implements the NodeChecker interface.
func (c *UnneededImportAliasChecker) Examples() []Example {
	return []Example{
		{
			Good: `
					import (
						"bar"
						foobar "foo/bar"
					)
				`,
			Bad: `
					import (
						"foo"
						qux "bar"
					)
				`,
		},
	}
}

// Register implements the NodeChecker interface.
func (c *UnneededImportAliasChecker) Register(fc *FileChecker) {
	fc.On(&ast.File{}, c)
}

// Check implements the NodeChecker interface.
func (c *UnneededImportAliasChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	file := node.(*ast.File)

	packageNames := map[string]struct{}{}
	for _, importSpec := range file.Imports {
		if hasAlias(importSpec) {
			continue
		}

		packageNames[extractPackageName(importSpec)] = struct{}{}
	}

	for _, importSpec := range file.Imports {
		if !hasAlias(importSpec) {
			continue
		}

		aliasName := importSpec.Name.Name
		if aliasName == "_" || aliasName == "." {
			continue
		}

		packageName := extractPackageName(importSpec)
		if _, ok := packageNames[packageName]; ok {
			continue
		}

		report.Errors = append(report.Errors, Error{
			Pos:     importSpec.Pos(),
			Message: fmt.Sprintf("unneeded package alias: %s", aliasName),
		})
	}
}

func hasAlias(importSpec *ast.ImportSpec) bool {
	return importSpec.Name != nil
}

func extractPackageName(importSpec *ast.ImportSpec) string {
	return path.Base(strings.Trim(importSpec.Path.Value, `"`))
}
