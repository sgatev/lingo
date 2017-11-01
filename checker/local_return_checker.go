package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register("local_return", NewLocalReturnChecker))
}

// LocalReturnChecker checks that exported funcs return exported
// (and internal) types only.
type LocalReturnChecker struct{}

// NewLocalReturnChecker constructs a LocalReturnChecker.
func NewLocalReturnChecker(configData interface{}) NodeChecker {
	return &LocalReturnChecker{}
}

// Register implements the NodeChecker interface.
func (c *LocalReturnChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncDecl{}, c)
	fc.On(&ast.Field{}, c)
}

// Check implements the NodeChecker interface.
func (c *LocalReturnChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	switch node := node.(type) {
	case *ast.FuncDecl:
		c.checkFuncDecl(node, report)
	case *ast.Field:
		c.checkField(node, report)
	}
}

func (c *LocalReturnChecker) checkFuncDecl(decl *ast.FuncDecl, report *Report) {
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

func (c *LocalReturnChecker) checkField(field *ast.Field, report *Report) {
	typ, ok := field.Type.(*ast.FuncType)
	if !ok {
		return
	}

	if typ.Results == nil {
		return
	}

	for _, name := range field.Names {
		if !name.IsExported() {
			continue
		}

		for _, result := range typ.Results.List {
			c.checkExpr(name.Name, result.Type, report)
		}
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

const localReturnErrMsg = "exported func '%s' cannot return value " +
	"of local type '%s'"

var internalTypes = map[string]struct{}{
	"int": {}, "int8": {},
	"int16": {}, "int32": {},
	"int64": {}, "uint": {},
	"uint8": {}, "uint16": {},
	"uint32": {}, "uint64": {},
	"byte": {}, "string": {},
	"float32": {}, "float64": {},
	"complex64": {}, "complex128": {},
	"bool": {}, "error": {},
}
