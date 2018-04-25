package checker

import (
	"fmt"
	"go/ast"

	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("func_results_count", NewFuncResultsCountChecker))
}

// FuncResultsCountConfig describes the configuration of a FuncResultsCountChecker.
type FuncResultsCountConfig struct {

	// Max is the maximum number of results of a func.
	Max int `mapdecode:"max"`
}

// FuncResultsCountChecker checks that funcs have a limited number of results.
type FuncResultsCountChecker struct {
	max int
}

// NewFuncResultsCountChecker constructs a FuncResultsCountChecker.
func NewFuncResultsCountChecker(configData interface{}) NodeChecker {
	var config FuncResultsCountConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &FuncResultsCountChecker{
		max: config.Max,
	}
}

// Title implements the NodeChecker interface.
func (c *FuncResultsCountChecker) Title() string {
	return "Func Results Count"
}

// Description implements the NodeChecker interface.
func (c *FuncResultsCountChecker) Description() string {
	return fmt.Sprintf(`The maximum number of results of a func is %d.`, c.max)
}

// Examples implements the NodeChecker interface.
func (c *FuncResultsCountChecker) Examples() []Example {
	return nil
}

// Register implements the NodeChecker interface.
func (c *FuncResultsCountChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncType{}, c)
}

// Check implements the NodeChecker interface.
func (c *FuncResultsCountChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	funcType := node.(*ast.FuncType)

	if funcType.Results == nil {
		return
	}

	var resultsCount int
	for _, result := range funcType.Results.List {
		if len(result.Names) > 0 {
			resultsCount += len(result.Names)
		} else {
			resultsCount += 1
		}
	}

	if resultsCount > c.max {
		report.Errors = append(report.Errors, Error{
			Pos: funcType.Pos(),
			Message: fmt.Sprintf("func has %d results, max is %d",
				resultsCount, c.max),
		})
	}
}
