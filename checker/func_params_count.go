package checker

import (
	"fmt"
	"go/ast"

	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("func_params_count", NewFuncParamsCountChecker))
}

// FuncParamsCountConfig describes the configuration of a FuncParamsCountChecker.
type FuncParamsCountConfig struct {

	// Max is the maximum number of parameters of a func.
	Max int `mapdecode:"max"`
}

// FuncParamsCountChecker checks that funcs have a limited number of parameters.
type FuncParamsCountChecker struct {
	max int
}

// NewFuncParamsCountChecker constructs a FuncParamsCountChecker.
func NewFuncParamsCountChecker(configData interface{}) NodeChecker {
	var config FuncParamsCountConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &FuncParamsCountChecker{
		max: config.Max,
	}
}

// Title implements the NodeChecker interface.
func (c *FuncParamsCountChecker) Title() string {
	return "Func Parameters Count"
}

// Description implements the NodeChecker interface.
func (c *FuncParamsCountChecker) Description() string {
	return fmt.Sprintf(`The maximum number of parameters of a func is %d.`, c.max)
}

// Examples implements the NodeChecker interface.
func (c *FuncParamsCountChecker) Examples() []Example {
	return nil
}

// Register implements the NodeChecker interface.
func (c *FuncParamsCountChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncType{}, c)
}

// Check implements the NodeChecker interface.
func (c *FuncParamsCountChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	funcType := node.(*ast.FuncType)

	var paramsCount int
	for _, param := range funcType.Params.List {
		if len(param.Names) > 0 {
			paramsCount += len(param.Names)
		} else {
			paramsCount += 1
		}
	}

	if paramsCount > c.max {
		report.Errors = append(report.Errors, Error{
			Pos: funcType.Pos(),
			Message: fmt.Sprintf("func has %d params, max is %d",
				paramsCount, c.max),
		})
	}
}
