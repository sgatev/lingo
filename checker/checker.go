package checker

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Error is a description of a checker violation.
type Error struct {

	// Pos is the position in a file where the error occurred.
	Pos token.Pos

	// Message is the error message.
	Message string
}

// Report collects the results of a run of some checkers.
type Report struct {

	// Errors contains all violations registered by the checkers.
	Errors []Error
}

// Example shows how to adhere and not adere to a rule.
type Example struct {

	// Good is an example of sticking to the rule.
	Good string

	// Bad is a counter-example showing a mis-use which lingo will report.
	Bad string
}

// NodeChecker checks ast.Node values for violations.
type NodeChecker interface {

	// Title returns the title of the node checker.
	Title() string

	// Description returns the detailed description of the node checker.
	Description() string

	// Examples is a set of examples that demonstrate the node checker rule.
	Examples() []Example

	// Register registers the node checker for specific types
	// of nodes in `fc`.
	Register(fc *FileChecker)

	// Check checks `node` and registers violations in `report`.
	Check(node ast.Node, content string, report *Report)
}

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

// Register registers all `checkers`.
func (c *FileChecker) Register(checkers ...NodeChecker) {
	for _, checker := range checkers {
		checker.Register(c)
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
func (c *FileChecker) Check(file *ast.File, content string, report *Report) {
	c.emit(file, content, report)

	for _, decl := range file.Decls {
		c.visitDecl(decl, report)
	}
}

func (c *FileChecker) emit(node ast.Node, content string, report *Report) {
	typeName := reflect.TypeOf(node).String()

	for _, checker := range c.checkers[typeName] {
		emittedContent := ""
		if len(content) > 0 {
			emittedContent = content[node.Pos()-1 : node.End()-1]
		}
		checker.Check(node, emittedContent, report)
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
	c.emit(decl, "", report)

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
	c.emit(spec, "", report)

	c.visitIdent(spec.Name, report)
	c.visitExpr(spec.Type, report)
}

func (c *FileChecker) visitValueSpec(spec *ast.ValueSpec, report *Report) {
	for _, name := range spec.Names {
		c.visitIdent(name, report)
	}
}

func (c *FileChecker) visitFuncDecl(decl *ast.FuncDecl, report *Report) {
	c.emit(decl, "", report)

	c.visitIdent(decl.Name, report)
	c.visitFuncType(decl.Type, report)
	c.visitBlockStmt(decl.Body, report)
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

func (c *FileChecker) visitRangeStmt(stmt *ast.RangeStmt, report *Report) {
	c.emit(stmt, "", report)

	c.visitBlockStmt(stmt.Body, report)
}

func (c *FileChecker) visitExpr(expr ast.Expr, report *Report) {
	switch expr := expr.(type) {
	case *ast.Ident:
		c.visitIdent(expr, report)
	case *ast.BinaryExpr:
		c.visitBinaryExpr(expr, report)
	case *ast.FuncLit:
		c.visitFuncLit(expr, report)
	case *ast.StructType:
		c.visitStructType(expr, report)
	case *ast.InterfaceType:
		c.visitInterfaceType(expr, report)
	case *ast.FuncType:
		c.visitFuncType(expr, report)
	case *ast.CallExpr:
		c.visitCallExpr(expr, report)
	}
}

func (c *FileChecker) visitCallExpr(expr *ast.CallExpr, report *Report) {
	c.visitExpr(expr.Fun, report)
}

func (c *FileChecker) visitIdent(ident *ast.Ident, report *Report) {
	c.emit(ident, "", report)
}

func (c *FileChecker) visitBinaryExpr(expr *ast.BinaryExpr, report *Report) {
	c.emit(expr, "", report)
	c.visitExpr(expr.X, report)
	c.visitExpr(expr.Y, report)
}

func (c *FileChecker) visitFuncLit(lit *ast.FuncLit, report *Report) {
	c.visitFuncType(lit.Type, report)
	c.visitBlockStmt(lit.Body, report)
}

func (c *FileChecker) visitBlockStmt(stmt *ast.BlockStmt, report *Report) {
	for _, stmt := range stmt.List {
		switch stmt := stmt.(type) {
		case *ast.DeclStmt:
			c.visitDeclStmt(stmt, report)
		case *ast.AssignStmt:
			c.visitAssignStmt(stmt, report)
		case *ast.RangeStmt:
			c.visitRangeStmt(stmt, report)
		case *ast.ExprStmt:
			c.visitExprStmt(stmt, report)
		case *ast.GoStmt:
			c.visitGoStmt(stmt, report)
		}
	}
}

func (c *FileChecker) visitGoStmt(stmt *ast.GoStmt, report *Report) {
	c.visitCallExpr(stmt.Call, report)
}

func (c *FileChecker) visitExprStmt(stmt *ast.ExprStmt, report *Report) {
	c.visitExpr(stmt.X, report)
}

func (c *FileChecker) visitStructType(typ *ast.StructType, report *Report) {
	c.emit(typ, "", report)

	for _, field := range typ.Fields.List {
		c.visitField(field, report)
	}
}

func (c *FileChecker) visitInterfaceType(
	typ *ast.InterfaceType,
	report *Report) {

	c.emit(typ, "", report)

	for _, method := range typ.Methods.List {
		c.visitField(method, report)
	}
}

func (c *FileChecker) visitFuncType(typ *ast.FuncType, report *Report) {
	c.emit(typ, "", report)
}

func (c *FileChecker) visitField(field *ast.Field, report *Report) {
	c.emit(field, "", report)

	for _, name := range field.Names {
		c.visitIdent(name, report)
	}

	c.visitExpr(field.Type, report)
}
