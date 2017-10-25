package checker_test

import (
	"fmt"
	"go/ast"
	"testing"

	. "github.com/s2gatev/lingo/checker"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	checker := &dummyChecker{}
	err := Register(func() NodeChecker {
		return checker
	})
	assert.Nil(t, err)
	assert.Equal(t, checker, Get(checker.Slug()))
}

func TestRegistryAlreadyPresent(t *testing.T) {
	checker := &dummyChecker{}
	err := Register(func() NodeChecker {
		return checker
	})
	assert.Equal(t, fmt.Errorf("checker already registered: dummy"), err)
	assert.Equal(t, checker, Get(checker.Slug()))
}

type dummyChecker struct{}

func (c *dummyChecker) Slug() string {
	return "dummy"
}

func (c *dummyChecker) Register(fc *FileChecker) {}

func (c *dummyChecker) Check(node ast.Node, report *Report) {}
