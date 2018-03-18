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
	assessed map[token.Pos]struct{}
}

// NewLeftQuantifiersChecker constructs a LeftQuantifiersChecker.
func NewLeftQuantifiersChecker(configData interface{}) NodeChecker {
	return &LeftQuantifiersChecker{
		assessed: map[token.Pos]struct{}{},
	}
}

// Title implements the NodeChecker interface.
func (c *LeftQuantifiersChecker) Title() string {
	return "Left Expression Quantifiers"
}

// Description implements the NodeChecker interface.
func (c *LeftQuantifiersChecker) Description() string {
	return `When a number literal appears in a binary expression it must be ` +
		`the left operand.`
}

// Examples implements the NodeChecker interface.
func (c *LeftQuantifiersChecker) Examples() []Example {
	return []Example{
		{
			Good: `_ = 5 * time.Minute`,
			Bad:  `_ = time.Minute * 5`,
		},
	}
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
	if _, ok := c.assessed[expr.OpPos]; ok {
		return
	}

	assessment := c.assess(node)
	if _, ok := validAssessments[assessment]; !ok {
		report.Errors = append(report.Errors, Error{
			Pos:     node.Pos(),
			Message: fmt.Sprintf("the left operand should be a basic literal"),
		})
	}
}

// assess determines the type of the given expression.
func (c *LeftQuantifiersChecker) assess(node ast.Node) assessment {
	switch expr := node.(type) {
	case *ast.BinaryExpr:
		c.assessed[expr.OpPos] = struct{}{}
		return c.assessBinaryExpr(expr)

	case *ast.BasicLit:
		return allQuantifiers

	case *ast.ParenExpr:
		return c.assess(expr.X)

	default:
		return noQuantifiers
	}
}

// assessBinaryExpr determines the type of the given binary expression.
func (c *LeftQuantifiersChecker) assessBinaryExpr(
	expr *ast.BinaryExpr) assessment {

	if _, ok := commutativeOperators[expr.Op]; !ok {
		return nonCommutative
	}

	x := c.assess(expr.X)
	y := c.assess(expr.Y)

	switch {
	case x == allQuantifiers && y == allQuantifiers:
		// e.g. 3 * 4 & 5
		return allQuantifiers

	case x == allQuantifiers && y == noQuantifiers:
		// e.g. 3 * 4 & a
		return leftQuantifiers

	case x == allQuantifiers && y == leftQuantifiers:
		// e.g. 3 * 4 & 5 * a
		return leftQuantifiers

	case x == leftQuantifiers && y == noQuantifiers:
		// e.g. 3 * a & a
		return leftQuantifiers

	case x == noQuantifiers && y == noQuantifiers:
		// e.g. a & a
		return noQuantifiers

	case x == nonCommutative || y == nonCommutative:
		// e.g. n - 1
		return nonCommutative

	default:
		return mixedQuantifiers
	}
}

// assessment indicates the type of an expression.
type assessment uint

const (
	// allQuantifiers indicates that an expression contains only basic literals,
	// e.g. 2 * 3, including single basic literals, e.g. 100.
	allQuantifiers assessment = 0

	// leftQuantifiers indicates that all quantifiers of an expression are on
	// the left side, e.g. 5 * a.
	leftQuantifiers assessment = 1

	// mixedQuantifiers indicates that the quantifiers of an expression are not
	// all on the left side, e.g. a * 5.
	mixedQuantifiers assessment = 2

	// noQuantifiers indicates that an expression contains no quantifiers,
	// e.g. a * a, including a single non-basic literal, e.g. x.
	noQuantifiers assessment = 3

	// notCommutative indicates that an expression contains a non-commutative
	// operator.
	nonCommutative assessment = 4
)

var validAssessments = map[assessment]struct{}{
	allQuantifiers:  struct{}{},
	leftQuantifiers: struct{}{},
	noQuantifiers:   struct{}{},
	nonCommutative:  struct{}{},
}

var commutativeOperators = map[token.Token]struct{}{
	token.ADD:  struct{}{}, // +
	token.MUL:  struct{}{}, // *
	token.AND:  struct{}{}, // &
	token.OR:   struct{}{}, // |
	token.XOR:  struct{}{}, // ^
	token.LAND: struct{}{}, // &&
	token.LOR:  struct{}{}, // ||
}
