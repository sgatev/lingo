package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register(&ConsistentReceiverNamesChecker{}))
}

// ConsistentReceiverNamesChecker checks that method receivers of a type
// are named consistently.
type ConsistentReceiverNamesChecker struct {
	receiverNames map[string]string
}

// Slug implements the NodeChecker interface.
func (c *ConsistentReceiverNamesChecker) Slug() string {
	return "consistent_receiver_names"
}

// Register implements the NodeChecker interface.
func (c *ConsistentReceiverNamesChecker) Register(fc *FileChecker) {
	fc.On(&ast.File{}, c)
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *ConsistentReceiverNamesChecker) Check(node ast.Node, report *Report) {
	switch node := node.(type) {
	case *ast.File:
		c.checkFile(node, report)
	case *ast.FuncDecl:
		c.checkFuncDecl(node, report)
	}
}

func (c *ConsistentReceiverNamesChecker) checkFile(
	file *ast.File,
	report *Report) {

	c.receiverNames = map[string]string{}
}

func (c *ConsistentReceiverNamesChecker) checkFuncDecl(
	decl *ast.FuncDecl,
	report *Report) {

	if decl.Recv == nil || len(decl.Recv.List) == 0 {
		return
	}

	names := decl.Recv.List[0].Names
	if len(names) == 0 {
		return
	}
	name := names[0].Name

	var typeName string
	switch typ := decl.Recv.List[0].Type.(type) {
	case *ast.Ident:
		typeName = typ.Name
	case *ast.StarExpr:
		if id, ok := typ.X.(*ast.Ident); ok {
			typeName = id.Name
		}
	}

	expectedName, ok := c.receiverNames[typeName]
	if !ok {
		c.receiverNames[typeName] = name
	} else if name != expectedName {
		report.Errors = append(report.Errors,
			fmt.Errorf("receivers in methods for type '%s' should have the same names",
				typeName))
	}
}
