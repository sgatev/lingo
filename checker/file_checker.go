package checker

import (
	"go/ast"
	"reflect"
)

// FileChecker checks ast.File values for violations.
type FileChecker struct {
	checkers map[string][]NodeChecker
}

// NewFileChecker creates a new FileChecker.
func NewFileChecker() *FileChecker {
	return &FileChecker{
		checkers: map[string][]NodeChecker{},
	}
}

// On registers `checker` for specific node type inferred from the type
// of `nodeType`.
func (c *FileChecker) On(nodeType interface{}, checker NodeChecker) {
	typeName := reflect.TypeOf(nodeType).String()

	c.checkers[typeName] = append(
		c.checkers[typeName], checker)
}

// Check checks `file` for violations and registers them in `report`.
func (c *FileChecker) Check(file *ast.File, report *Report) {
	for _, decl := range file.Decls {
		c.visitDecl(decl, report)
	}
}

func (c *FileChecker) visitDecl(decl ast.Decl, report *Report) {
	switch decl := decl.(type) {
	case *ast.GenDecl:
		c.visitGenDecl(decl, report)
	case *ast.FuncDecl:
		c.visitFuncDecl(decl, report)
	}
}

func (c *FileChecker) visitGenDecl(decl *ast.GenDecl, report *Report) {
	for _, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.TypeSpec:
			c.visitTypeSpec(spec, report)
		}
	}
}

func (c *FileChecker) visitTypeSpec(spec *ast.TypeSpec, report *Report) {
	typeName := reflect.TypeOf(spec).String()
	for _, checker := range c.checkers[typeName] {
		checker.Check(spec, report)
	}
}

func (c *FileChecker) visitFuncDecl(decl *ast.FuncDecl, report *Report) {
	for _, stmt := range decl.Body.List {
		switch stmt := stmt.(type) {
		case *ast.DeclStmt:
			c.visitDeclStmt(stmt, report)
		}
	}
}

func (c *FileChecker) visitDeclStmt(stmt *ast.DeclStmt, report *Report) {
	c.visitDecl(stmt.Decl, report)
}
