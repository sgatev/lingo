package checker

import (
	"fmt"
	"go/ast"
)

// ExportedIdentDocChecker checks the documentation of exported
// identifiers.
type ExportedIdentDocChecker struct{}

// Register implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Register(fc *FileChecker) {
	fc.On(&ast.GenDecl{}, c)
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Check(node ast.Node, report *Report) {
	switch decl := node.(type) {
	case *ast.GenDecl:
		for _, spec := range decl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				c.checkTypeSpec(decl.Doc, spec, report)
			case *ast.ValueSpec:
				c.checkValueSpec(decl.Doc, spec, report)
			}
		}
	case *ast.FuncDecl:
		c.checkFuncDecl(decl, report)
	}
}

func (c *ExportedIdentDocChecker) checkTypeSpec(
	doc *ast.CommentGroup,
	spec *ast.TypeSpec,
	report *Report) {

	if !spec.Name.IsExported() {
		return
	}

	if doc == nil {
		report.Errors = append(report.Errors,
			fmt.Errorf("exported identifier '%s' is not documented",
				spec.Name.Name))
		return
	}
}

func (c *ExportedIdentDocChecker) checkValueSpec(
	doc *ast.CommentGroup,
	spec *ast.ValueSpec,
	report *Report) {

	for _, name := range spec.Names {
		if !name.IsExported() {
			continue
		}

		if doc == nil {
			report.Errors = append(report.Errors,
				fmt.Errorf("exported identifier '%s' is not documented",
					name.Name))
			continue
		}
	}
}

func (c *ExportedIdentDocChecker) checkFuncDecl(
	decl *ast.FuncDecl,
	report *Report) {

	if !decl.Name.IsExported() {
		return
	}

	if decl.Doc == nil {
		report.Errors = append(report.Errors,
			fmt.Errorf("exported identifier '%s' is not documented",
				decl.Name.Name))
		return
	}
}
