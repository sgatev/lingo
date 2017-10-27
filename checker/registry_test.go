package checker_test

import (
	"fmt"
	"go/ast"
	"testing"

	. "github.com/s2gatev/lingo/checker"
	"github.com/stretchr/testify/assert"
)

func TestRegistryRegister(t *testing.T) {
	checker := &dummyChecker{}
	err := Register(func() NodeChecker {
		return checker
	})
	assert.Nil(t, err)
	assert.Equal(t, checker, Get(checker.Slug()))
}

func TestRegistryRegisterAlreadyPresent(t *testing.T) {
	checker := &dummyChecker{}
	err := Register(func() NodeChecker {
		return checker
	})
	assert.Equal(t, fmt.Errorf("checker already registered: dummy"), err)
	assert.Equal(t, checker, Get(checker.Slug()))
}

func TestRegistryGetNotPresent(t *testing.T) {
	assert.Nil(t, Get("unknown"))
}

type dummyChecker struct{}

func (c *dummyChecker) Slug() string {
	return "dummy"
}

func (c *dummyChecker) Register(fc *FileChecker) {}

func (c *dummyChecker) Check(node ast.Node, report *Report) {}
