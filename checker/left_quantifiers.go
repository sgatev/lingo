package checker

import (
	"fmt"
	"go/ast"
	"go/token"
)

func init() {
	must(Register("left_quantifiers", NewLeftQuantifiersChecker))
}

// LeftQuantifiersChecker checks that when a basic literal appears in a binary
// expression it is the left operand.
type LeftQuantifiersChecker struct {
	errorPositions map[token.Pos]struct{}
}

// NewLeftQuantifiersChecker constructs a BasicLiteralLeftOperandChecker.
func NewLeftQuantifiersChecker(configData interface{}) NodeChecker {
	return &LeftQuantifiersChecker{}
}

// Register implements the NodeChecker interface.
func (c *LeftQuantifiersChecker) Register(fc *FileChecker) {
	fc.On(&ast.BinaryExpr{}, c)
}

// Check implements the NodeChecker interface.
func (c *LeftQuantifiersChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	expr := node.(*ast.BinaryExpr)
	if _, ok := c.errorPositions[expr.OpPos]; ok {
		return
	}

	positions := map[token.Pos]struct{}{}
	if !assertLeftQuantifiers(node, positions) {
		c.errorPositions = positions
		report.Errors = append(report.Errors,
			fmt.Errorf("the left operand should be a basic literal"))
	}

}

func assertLeftQuantifiers(node ast.Node, pos map[token.Pos]struct{}) bool {
	switch expr := node.(type) {
	case *ast.BinaryExpr:
		pos[expr.OpPos] = struct{}{}
		if assertBasicLit(expr.X, pos) {
			return assertLeftQuantifiers(expr.Y, pos)
		}
		return !assertBasicLit(expr.Y, pos)
	case *ast.ParenExpr:
		return assertLeftQuantifiers(expr.X, pos)
	default:
		return true
	}
}

func assertBasicLit(node ast.Node, pos map[token.Pos]struct{}) bool {
	switch expr := node.(type) {
	case *ast.BinaryExpr:
		pos[expr.OpPos] = struct{}{}
		return assertBasicLit(expr.X, pos) && assertBasicLit(expr.Y, pos)
	case *ast.ParenExpr:
		return assertBasicLit(expr.X, pos)
	case *ast.BasicLit:
		return true
	default:
		return false
	}
}
