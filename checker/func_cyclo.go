package checker

import (
	"fmt"
	"go/ast"

	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("func_cyclo", NewFuncCycloChecker))
}

// FuncCycloConfig describes the configuration of a FuncCycloChecker.
type FuncCycloConfig struct {

	// Max is the maximum cyclomatic complexity of a func.
	Max int `mapdecode:"max"`
}

// FuncCycloChecker checks that funcs are within specific cyclomatic complexity.
type FuncCycloChecker struct {
	max int
}

// NewFuncCycloChecker constructs a FuncCycloChecker.
func NewFuncCycloChecker(configData interface{}) NodeChecker {
	var config FuncCycloConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &FuncCycloChecker{
		max: config.Max,
	}
}

// Title implements the NodeChecker interface.
func (c *FuncCycloChecker) Title() string {
	return "Func Cyclo Complexity"
}

// Description implements the NodeChecker interface.
func (c *FuncCycloChecker) Description() string {
	// TODO: add reference to cyclomatic complexity definition
	return fmt.Sprintf(`The maximum cyclomatic complexity of a func is %d.`, c.max)
}

// Examples implements the NodeChecker interface.
func (c *FuncCycloChecker) Examples() []Example {
	return nil
}

// Register implements the NodeChecker interface.
func (c *FuncCycloChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *FuncCycloChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	funcDecl := node.(*ast.FuncDecl)
	var complexity complexityComputer
	ast.Walk(&complexity, funcDecl)
	if int(complexity) > c.max {
		report.Errors = append(report.Errors, Error{
			Pos: funcDecl.Pos(),
			Message: fmt.Sprintf("func %s has cyclomatic complexity %d, max is %d",
				funcDecl.Name.Name, complexity, c.max),
		})
	}
}

type complexityComputer int

// Visit implements the ast.Visitor interface.
func (c *complexityComputer) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.FuncDecl,
		*ast.IfStmt,
		*ast.ForStmt,
		*ast.RangeStmt,
		*ast.CaseClause,
		*ast.CommClause:

		*c++
	}

	return c
}
