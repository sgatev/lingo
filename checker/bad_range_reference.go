package checker

import (
	"fmt"
	"go/ast"
	"go/token"
)

func init() {
	must(Register("bad_range_reference", NewBadRangeReferenceChecker))
}

// BadRangeReferenceChecker checks that vars declared in a range statement
// are not used by reference inside the body of the loop.
type BadRangeReferenceChecker struct{}

// NewBadRangeReferenceChecker constructs a new BadRangeReferenceChecker.
func NewBadRangeReferenceChecker(configData interface{}) NodeChecker {
	return &BadRangeReferenceChecker{}
}

// Title implements the NodeChecker interface.
func (c *BadRangeReferenceChecker) Title() string {
	return "Bad Range Reference"
}

// Description implements the NodeChecker interface.
func (c *BadRangeReferenceChecker) Description() string {
	return `A value declared in a range statement must not be used by reference.`
}

// Examples implements the NodeChecker interface.
func (c *BadRangeReferenceChecker) Examples() []Example {
	return []Example{
		{
			Good: `
			for _, value := range values {
				foo(value)
			}
			`,
			Bad: `
			for _, value := range values {
				foo(&value)
			}
			`,
		},
		{
			Good: `
			for _, value := range values {
				value := value
				foo(&value)
			}
			`,
			Bad: `
			for _, value := range values {
				value := &value
				foo(value)
			}
			`,
		},
	}
}

// Register implements the NodeChecker interface.
func (c *BadRangeReferenceChecker) Register(fc *FileChecker) {
	fc.On(&ast.RangeStmt{}, c)
}

// Check implements the NodeChecker interface.
func (c *BadRangeReferenceChecker) Check(
	node ast.Node, content string, report *Report) {

	stmt := node.(*ast.RangeStmt)
	c.checkReference(report, stmt.Key, stmt.Body)
	c.checkReference(report, stmt.Value, stmt.Body)
}

func (c *BadRangeReferenceChecker) checkReference(
	report *Report, expr ast.Expr, body *ast.BlockStmt) {

	if expr == nil {
		return
	}

	ident := expr.(*ast.Ident)
	if ident.Obj.Name == "_" {
		return
	}

	ast.Walk(&badReferenceVisitor{
		report: report,
		name:   ident.Obj.Name,
	}, body)
}

type badReferenceVisitor struct {
	report     *Report
	name       string
	overridden bool
}

// Visit implements the ast.Visitor interface.
func (v *badReferenceVisitor) Visit(node ast.Node) ast.Visitor {
	if v.overridden {
		return nil
	}

	switch node := node.(type) {
	case *ast.CallExpr:
		v.visitCallExpr(node)
	case *ast.AssignStmt:
		v.visitAssignStmt(node)
	}

	return v
}

func (v *badReferenceVisitor) visitCallExpr(expr *ast.CallExpr) {
	switch expr.Fun.(type) {
	case *ast.FuncLit:
		v.visitAnonymousFuncCallExpr(expr.Fun.(*ast.FuncLit))
	case *ast.Ident:
		v.visitFuncCallExpr(expr)
	case *ast.SelectorExpr:
		v.visitFuncCallExpr(expr)
	}
}

func (v *badReferenceVisitor) visitAnonymousFuncCallExpr(expr *ast.FuncLit) {
	for _, field := range expr.Type.Params.List {
		for _, name := range field.Names {
			if name.Name == v.name {
				v.overridden = true
			}
		}
	}
}

func (v *badReferenceVisitor) visitFuncCallExpr(expr *ast.CallExpr) {
	for _, arg := range expr.Args {
		switch expr := arg.(type) {
		case *ast.UnaryExpr:
			v.checkUnaryExpr(expr)
		}
	}
}

func (v *badReferenceVisitor) visitAssignStmt(stmt *ast.AssignStmt) {
	for _, expr := range stmt.Rhs {
		expr, ok := expr.(*ast.UnaryExpr)
		if !ok {
			continue
		}

		v.checkUnaryExpr(expr)
	}
	for _, expr := range stmt.Lhs {
		ident, ok := expr.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Obj.Name == v.name {
			v.overridden = true
			break
		}
	}
}

func (v *badReferenceVisitor) checkUnaryExpr(expr *ast.UnaryExpr) {
	ident, ok := expr.X.(*ast.Ident)
	if !ok {
		return
	}

	if ident.Obj.Name != v.name {
		return
	}

	if expr.Op == token.AND {
		v.report.Errors = append(v.report.Errors, Error{
			Pos:     expr.Pos(),
			Message: fmt.Sprintf("bad reference of range var: %s", v.name),
		})
	}
}
