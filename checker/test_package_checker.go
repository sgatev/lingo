package checker

import (
	"fmt"
	"go/ast"
	"strings"
)

func init() {
	must(Register(NewTestPackageChecker))
}

// TestPackageChecker checks that tests are placed in "*_test" packages
// only.
type TestPackageChecker struct{}

// NewTestPackageChecker constructs a TestPackageChecker.
func NewTestPackageChecker() NodeChecker {
	return &TestPackageChecker{}
}

// Slug implements the NodeChecker interface.
func (c *TestPackageChecker) Slug() string {
	return "test_package"
}

// Register implements the NodeChecker interface.
func (c *TestPackageChecker) Register(fc *FileChecker) {
	fc.On(&ast.File{}, c)
}

// Check implements the NodeChecker interface.
func (c *TestPackageChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	file := node.(*ast.File)

	isTestPackage := false
	for _, importSpec := range file.Imports {
		if importSpec.Path.Value == `"testing"` {
			isTestPackage = true
			break
		}
	}

	if !isTestPackage {
		return
	}

	packageName := file.Name.Name
	if strings.HasSuffix(packageName, "_test") {
		return
	}

	report.Errors = append(report.Errors,
		fmt.Errorf("package '%s' should be named '%s_test'",
			packageName, packageName))
}
