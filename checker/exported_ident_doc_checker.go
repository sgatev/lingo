package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register(NewExportedIdentDocChecker))
}

// ExportedIdentDocChecker checks the documentation of exported
// identifiers.
type ExportedIdentDocChecker struct{}

// NewExportedIdentDocChecker constructs a ExportedIdentDocChecker.
func NewExportedIdentDocChecker() NodeChecker {
	return &ExportedIdentDocChecker{}
}

// Slug implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Slug() string {
	return "exported_ident_doc"
}

// Register implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Register(fc *FileChecker) {
	fc.On(&ast.GenDecl{}, c)
	fc.On(&ast.FuncDecl{}, c)
	fc.On(&ast.Field{}, c)
}

// Check implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	switch node := node.(type) {
	case *ast.GenDecl:
		for _, spec := range node.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				c.checkTypeSpec(node.Doc, spec, report)
			case *ast.ValueSpec:
				c.checkValueSpec(node.Doc, spec, report)
			}
		}
	case *ast.FuncDecl:
		c.checkFuncDecl(node, report)
	case *ast.Field:
		c.checkField(node, report)
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

		if doc == nil && spec.Doc == nil {
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

func (c *ExportedIdentDocChecker) checkField(
	field *ast.Field,
	report *Report) {

	var exported []string
	for _, name := range field.Names {
		if name.IsExported() {
			exported = append(exported, name.Name)
		}
	}

	if len(exported) == 0 {
		return
	}

	if field.Doc == nil {
		for _, name := range exported {
			report.Errors = append(report.Errors,
				fmt.Errorf("exported identifier '%s' is not documented",
					name))
		}
		return
	}
}
