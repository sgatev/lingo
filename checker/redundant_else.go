package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register("redundant_else", NewRedundantElseChecker))
}

// RedundantElseChecker checks that if the body of an 'if' statement ends with a
// terminating statement, there is no 'else' statement.
type RedundantElseChecker struct{}

// NewRedundantElseChecker constructs a RedundantElseChecker.
func NewRedundantElseChecker(configData interface{}) NodeChecker {
	return &RedundantElseChecker{}
}

// Title implements the NodeChecker interface.
func (c *RedundantElseChecker) Title() string {
	return "Redundant Else"
}

// Description implements the NodeChecker interface.
func (c *RedundantElseChecker) Description() string {
	return `When an if statement ends with a terminating statement ` +
		`it should not be followed by an else statement.`
}

// Examples implements the NodeChecker interface.
func (c *RedundantElseChecker) Examples() []Example {
	return []Example{
		{
			Good: `
			if err != nil {
				return err
			}
			call(foo)
			`,
			Bad: `
			if err != nil {
				return err
			} else {
				call(foo)
			}
			`,
		},
	}
}

// Register implements the NodeChecker interface.
func (c *RedundantElseChecker) Register(fc *FileChecker) {
	fc.On(&ast.IfStmt{}, c)
}

// Check implements the NodeChecker interface.
func (c *RedundantElseChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	stmt := node.(*ast.IfStmt)
	if stmt.Else == nil {
		return
	}

	var termStmt string
	lastStmt := stmt.Body.List[len(stmt.Body.List)-1]
	switch lastStmt := lastStmt.(type) {
	case *ast.ReturnStmt:
		termStmt = "return"
	case *ast.BranchStmt:
		termStmt = lastStmt.Tok.String()
	case *ast.ExprStmt:
		expr, ok := lastStmt.X.(*ast.CallExpr)
		if !ok {
			break
		}

		if c.isPanic(expr) {
			termStmt = "panic()"
		} else if c.isExit(expr) {
			termStmt = "os.Exit()"
		}
	}

	if termStmt != "" {
		report.Errors = append(report.Errors, Error{
			Pos:     stmt.Else.Pos(),
			Message: fmt.Sprintf("unexpected else after %s statement", termStmt),
		})
	}
}

func (c *RedundantElseChecker) isPanic(expr *ast.CallExpr) bool {
	ident, ok := expr.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "panic"
}

func (c *RedundantElseChecker) isExit(expr *ast.CallExpr) bool {
	sel, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	// TODO: Handle the case where the os package is aliased in
	// the scope of the file.
	return ident.Name == "os" && sel.Sel.Name == "Exit"
}
