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

func (c *FileChecker) emit(node ast.Node, report *Report) {
	typeName := reflect.TypeOf(node).String()

	for _, checker := range c.checkers[typeName] {
		checker.Check(node, report)
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
	c.emit(decl, report)

	for _, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.TypeSpec:
			c.visitTypeSpec(spec, report)
		case *ast.ValueSpec:
			c.visitValueSpec(spec, report)
		}
	}
}

func (c *FileChecker) visitTypeSpec(spec *ast.TypeSpec, report *Report) {
	c.emit(spec, report)

	c.visitIdent(spec.Name, report)
	c.visitExpr(spec.Type, report)
}

func (c *FileChecker) visitValueSpec(spec *ast.ValueSpec, report *Report) {
	for _, name := range spec.Names {
		c.visitIdent(name, report)
	}
}

func (c *FileChecker) visitFuncDecl(decl *ast.FuncDecl, report *Report) {
	c.emit(decl, report)

	c.visitIdent(decl.Name, report)

	for _, stmt := range decl.Body.List {
		switch stmt := stmt.(type) {
		case *ast.DeclStmt:
			c.visitDeclStmt(stmt, report)
		case *ast.AssignStmt:
			c.visitAssignStmt(stmt, report)
		}
	}
}

func (c *FileChecker) visitDeclStmt(stmt *ast.DeclStmt, report *Report) {
	c.visitDecl(stmt.Decl, report)
}

func (c *FileChecker) visitAssignStmt(stmt *ast.AssignStmt, report *Report) {
	for _, expr := range stmt.Lhs {
		c.visitExpr(expr, report)
	}

	for _, expr := range stmt.Rhs {
		c.visitExpr(expr, report)
	}
}

func (c *FileChecker) visitExpr(expr ast.Expr, report *Report) {
	switch expr := expr.(type) {
	case *ast.Ident:
		c.visitIdent(expr, report)
	case *ast.FuncLit:
		c.visitFuncLit(expr, report)
	case *ast.StructType:
		c.visitStructType(expr, report)
	}
}

func (c *FileChecker) visitIdent(ident *ast.Ident, report *Report) {
	c.emit(ident, report)
}

func (c *FileChecker) visitFuncLit(lit *ast.FuncLit, report *Report) {
	for _, stmt := range lit.Body.List {
		switch stmt := stmt.(type) {
		case *ast.DeclStmt:
			c.visitDeclStmt(stmt, report)
		case *ast.AssignStmt:
			c.visitAssignStmt(stmt, report)
		}
	}
}

func (c *FileChecker) visitStructType(typ *ast.StructType, report *Report) {
	c.emit(typ, report)

	for _, field := range typ.Fields.List {
		c.visitField(field, report)
	}
}

func (c *FileChecker) visitField(field *ast.Field, report *Report) {
	c.emit(field, report)

	for _, name := range field.Names {
		c.visitIdent(name, report)
	}
}
