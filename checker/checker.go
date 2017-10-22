package checker

import "go/ast"

// Report collects the results of a run of some checkers.
type Report struct {

	// Errors contains all violations registered by the checkers.
	Errors []error
}

// NodeChecker checks ast.Node values for violations.
type NodeChecker interface {

	// Register registers the node checker for specific types
	// of nodes in `fc`.
	Register(fc *FileChecker)

	// Check checks `node` and registers violations in `report`.
	Check(node ast.Node, report *Report)
}
