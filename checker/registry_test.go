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
	err := Register("dummy", func(configData interface{}) NodeChecker {
		return checker
	})
	assert.Nil(t, err)
	assert.Equal(t, checker, Get("dummy", nil))
}

func TestRegistryRegisterAlreadyPresent(t *testing.T) {
	checker := &dummyChecker{}
	err := Register("dummy", func(configData interface{}) NodeChecker {
		return checker
	})
	assert.Equal(t, fmt.Errorf("checker already registered: dummy"), err)
	assert.Equal(t, checker, Get("dummy", nil))
}

func TestRegistryGetNotPresent(t *testing.T) {
	assert.Nil(t, Get("unknown", nil))
}

type dummyChecker struct{}

func (c *dummyChecker) Title() string {
	return ""
}

func (c *dummyChecker) Description() string {
	return ""
}

func (c *dummyChecker) Register(fc *FileChecker) {}

func (c *dummyChecker) Check(node ast.Node, content string, report *Report) {}
