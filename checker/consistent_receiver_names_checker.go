package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register(
		"consistent_receiver_names",
		NewConsistentReceiverNamesChecker))
}

// ConsistentReceiverNamesChecker checks that method receivers of a type
// are named consistently.
type ConsistentReceiverNamesChecker struct {
	receiverNames map[string]string
}

// NewConsistentReceiverNamesChecker constructs a
// ConsistentReceiverNamesChecker.
func NewConsistentReceiverNamesChecker(configData interface{}) NodeChecker {
	return &ConsistentReceiverNamesChecker{
		receiverNames: map[string]string{},
	}
}

// Register implements the NodeChecker interface.
func (c *ConsistentReceiverNamesChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *ConsistentReceiverNamesChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	decl := node.(*ast.FuncDecl)

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
			fmt.Errorf("receivers in methods for type '%s' "+
				"should have the same names",
				typeName))
	}
}
