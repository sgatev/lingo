package checker

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("exported_ident_doc", NewExportedIdentDocChecker))
}

// ExportedIdentDocCheckerConfig describes the configuration of a ExportedIdentDocChecker.
type ExportedIdentDocCheckerConfig struct {

	// HasIdentPrefix signals if the checker should ensure that every doc comment begins
	// with the name of the item it describes.
	HasIdentPrefix bool `mapdecode:"has_ident_prefix"`
}

// ExportedIdentDocChecker checks the documentation of exported
// identifiers.
type ExportedIdentDocChecker struct {
	hasIdentPrefix bool
}

// NewExportedIdentDocChecker constructs a ExportedIdentDocChecker.
func NewExportedIdentDocChecker(configData interface{}) NodeChecker {
	var config ExportedIdentDocCheckerConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		fmt.Println(err)
		return nil
	}

	return &ExportedIdentDocChecker{
		hasIdentPrefix: config.HasIdentPrefix,
	}
}

// Title implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Title() string {
	return "Documented Exported Identifiers"
}

// Description implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Description() string {
	description := `Every exported identifier must be documented.`
	if c.hasIdentPrefix {
		description += ` The documentation string must begin with the name of ` +
			`the identifier.`
	}
	return description
}

// Examples implements the NodeChecker interface.
func (c *ExportedIdentDocChecker) Examples() []Example {
	examples := []Example{
		{
			Good: `
// Runner can run and stop.
type Runner interface {

    // Run runs the runner.
    Run()

    // Stop stops the runner.
    Stop() error
}

// Runners are many runners.
var Runners []Runner
`,
			Bad: `
type Runner interface {

    Run()

    // This method is not documented properly.
    Stop() error
}

var Runners []Runner
`,
		},
	}
	return examples
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
		report.Errors = append(report.Errors, Error{
			Pos: spec.Pos(),
			Message: fmt.Sprintf("exported identifier '%s' is not documented",
				spec.Name.Name),
		})
		return
	}

	c.checkPrefix(spec.Name.Name, doc, report)
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
			report.Errors = append(report.Errors, Error{
				Pos: spec.Pos(),
				Message: fmt.Sprintf("exported identifier '%s' is not documented",
					name.Name),
			})
			continue
		}

		identDoc := spec.Doc
		if spec.Doc == nil {
			identDoc = doc
		}
		c.checkPrefix(name.Name, identDoc, report)
	}
}

func (c *ExportedIdentDocChecker) checkFuncDecl(
	decl *ast.FuncDecl,
	report *Report) {

	if !decl.Name.IsExported() {
		return
	}

	if decl.Doc == nil {
		report.Errors = append(report.Errors, Error{
			Pos: decl.Pos(),
			Message: fmt.Sprintf("exported identifier '%s' is not documented",
				decl.Name.Name),
		})
		return
	}

	c.checkPrefix(decl.Name.Name, decl.Doc, report)
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
			report.Errors = append(report.Errors, Error{
				Pos: field.Pos(),
				Message: fmt.Sprintf("exported identifier '%s' is not documented",
					name),
			})
		}
		return
	}

	if len(exported) == 1 {
		c.checkPrefix(exported[0], field.Doc, report)
	}
}

func (c *ExportedIdentDocChecker) checkPrefix(
	name string,
	doc *ast.CommentGroup,
	report *Report) {

	if !c.hasIdentPrefix {
		return
	}

	if !strings.HasPrefix(doc.Text(), name) {
		report.Errors = append(report.Errors, Error{
			Pos:     doc.Pos(),
			Message: fmt.Sprintf("expected the comment to start with '%s'", name),
		})
	}
}
