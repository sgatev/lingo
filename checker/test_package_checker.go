package checker

import (
	"fmt"
	"go/ast"
	"strings"
)

func init() {
	must(Register("test_package", NewTestPackageChecker))
}

// TestPackageChecker checks that tests are placed in "*_test" packages
// only.
type TestPackageChecker struct{}

// NewTestPackageChecker constructs a TestPackageChecker.
func NewTestPackageChecker(configData interface{}) NodeChecker {
	return &TestPackageChecker{}
}

// Title implements the NodeChecker interface.
func (c *TestPackageChecker) Title() string {
	return "Test Package"
}

// Description implements the NodeChecker interface.
func (c *TestPackageChecker) Description() string {
	return `Tests must be defined in a separate package.`
}

// Examples implements the NodeChecker interface.
func (c *TestPackageChecker) Examples() []Example {
	return []Example{
		{
			Good: `package feature_test

import "testing"

func TestFeature(t *testing.T) {}`,
			Bad: `package feature

import "testing"

func TestFeature(t *testing.T) {}`,
		},
	}
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

	report.Errors = append(report.Errors, Error{
		Pos: node.Pos(),
		Message: fmt.Sprintf("package '%s' should be named '%s_test'",
			packageName, packageName),
	})
}
