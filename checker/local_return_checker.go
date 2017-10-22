package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	Register(&LocalReturnChecker{})
}

// LocalReturnChecker checks that exported funcs return exported
// (and internal) types only.
type LocalReturnChecker struct{}

// Slug implements the NodeChecker interface.
func (c *LocalReturnChecker) Slug() string {
	return "local_return"
}

// Register implements the NodeChecker interface.
func (c *LocalReturnChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *LocalReturnChecker) Check(node ast.Node, report *Report) {
	decl := node.(*ast.FuncDecl)

	if !decl.Name.IsExported() {
		return
	}

	if decl.Type.Results == nil {
		return
	}

	for _, result := range decl.Type.Results.List {
		c.checkExpr(decl.Name.Name, result.Type, report)
	}
}

func (c *LocalReturnChecker) checkExpr(
	funcName string,
	expr ast.Expr,
	report *Report) {

	switch expr := expr.(type) {
	case *ast.Ident:
		c.checkIdent(funcName, expr, report)
	case *ast.ChanType:
		c.checkExpr(funcName, expr.Value, report)
	case *ast.ArrayType:
		c.checkExpr(funcName, expr.Elt, report)
	}
}

func (c *LocalReturnChecker) checkIdent(
	funcName string,
	ident *ast.Ident,
	report *Report) {

	if _, ok := internalTypes[ident.Name]; ok {
		return
	}

	if ident.IsExported() {
		return
	}

	report.Errors = append(report.Errors,
		fmt.Errorf(localReturnErrMsg,
			funcName, ident.Name))
}

const localReturnErrMsg = "exported func '%s' cannot return value of local type '%s'"

var internalTypes = map[string]struct{}{
	"int": struct{}{}, "int8": struct{}{},
	"int16": struct{}{}, "int32": struct{}{},
	"int64": struct{}{}, "uint": struct{}{},
	"uint8": struct{}{}, "uint16": struct{}{},
	"uint32": struct{}{}, "uint64": struct{}{},
	"byte": struct{}{}, "string": struct{}{},
	"float32": struct{}{}, "float64": struct{}{},
	"complex64": struct{}{}, "complex128": struct{}{},
	"bool": struct{}{}, "error": struct{}{},
}
