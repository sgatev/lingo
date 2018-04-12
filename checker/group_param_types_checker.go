package checker

import "go/ast"

func init() {
	must(Register("group_param_types", NewGroupParamTypesChecker))
}

// GroupParamTypesChecker checks that func parameters are grouped by
// type.
type GroupParamTypesChecker struct{}

// NewGroupParamTypesChecker constructs a GroupParamTypesChecker.
func NewGroupParamTypesChecker(configData interface{}) NodeChecker {
	return &GroupParamTypesChecker{}
}

// Title implements the NodeChecker interface.
func (c *GroupParamTypesChecker) Title() string {
	return "Group Param Types"
}

// Description implements the NodeChecker interface.
func (c *GroupParamTypesChecker) Description() string {
	return `Group parameters of the same type.`
}

// Examples implements the NodeChecker interface.
func (c *GroupParamTypesChecker) Examples() []Example {
	return []Example{
		{
			Good: `func foo(a, b string) {}`,
			Bad:  `func foo(a string, b string) {}`,
		},
	}
}

// Register implements the NodeChecker interface.
func (c *GroupParamTypesChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncType{}, c)
}

// Check implements the NodeChecker interface.
func (c *GroupParamTypesChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	funcType := node.(*ast.FuncType)

	var prevType string
	for _, param := range funcType.Params.List {
		curType := typeName(param.Type)
		if curType == prevType {
			report.Errors = append(report.Errors, Error{
				Pos:     node.Pos(),
				Message: `params should be grouped by type`,
			})
		}
		prevType = curType
	}
}

func typeName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StarExpr:
		return "*" + typeName(expr.X)
	case *ast.SelectorExpr:
		return typeName(expr.X) + "." + expr.Sel.Name
	case *ast.ChanType:
		var result string
		switch expr.Dir {
		case ast.SEND:
			result = "chan<-"
		case ast.RECV:
			result = "<-chan"
		default:
			result = "chan"
		}
		return result + " " + typeName(expr.Value)
	case *ast.Ellipsis:
		return "..." + typeName(expr.Elt)
	default:
		return ""
	}
}
