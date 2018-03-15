package checker

import (
	"fmt"
	"go/ast"
)

func init() {
	must(Register("pass_context_first", NewPassContextFirstChecker))
}

// PassContextFirstChecker checks that if a function declaration contains a context, it is
// the first argument.
type PassContextFirstChecker struct{}

// NewPassContextFirstChecker constructs a PassContextFirstChecker.
func NewPassContextFirstChecker(configData interface{}) NodeChecker {
	return &PassContextFirstChecker{}
}

// Title implements the NodeChecker interface.
func (c *PassContextFirstChecker) Title() string {
	return "Context Argument First"
}

// Description implements the NodeChecker interface.
func (c *PassContextFirstChecker) Description() string {
	return `A function must receive context.Context as its first argument.`
}

// Examples implements the NodeChecker interface.
func (c *PassContextFirstChecker) Examples() []Example {
	return []Example{
		{
			Good: `func Get(ctx context.Context, id string) {}`,
			Bad:  `func Get(id string, ctx context.Context) {}`,
		},
	}
}

// Register implements the NodeChecker interface.
func (c *PassContextFirstChecker) Register(fc *FileChecker) {
	fc.On(&ast.FuncDecl{}, c)
}

// Check implements the NodeChecker interface.
func (c *PassContextFirstChecker) Check(
	node ast.Node,
	content string,
	report *Report) {

	decl := node.(*ast.FuncDecl)

	if decl.Type.Params == nil {
		return
	}

	params := decl.Type.Params.List
	if len(params) == 0 {
		return
	}

	contextNotFirst := false
	for _, param := range params[1:] {
		selector, ok := param.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			continue
		}

		// TODO: Handle the case where the context package is aliased in the scope of
		// the file.
		if ident.Name == "context" && selector.Sel.Name == "Context" {
			contextNotFirst = true
		}
	}

	if contextNotFirst {
		report.Errors = append(report.Errors, Error{
			Pos: node.Pos(),
			Message: fmt.Sprintf("func '%s' should be passed context as first parameter",
				decl.Name.Name),
		})
	}
}
