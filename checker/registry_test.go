package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	checker := &LocalReturnChecker{}
	Register(checker)
	assert.Equal(t, checker, Get(checker.Slug()))
}
