package checker_test

import (
	"testing"

	. "github.com/s2gatev/lingo/checker"

	"github.com/stretchr/testify/assert"
)

func TestTypeNameCorrectExported(t *testing.T) {
	file := ParseFileContent(`
		package test

		type FooBar struct{}
	`)

	checker := &TypeNameChecker{}
	err := checker.Check(file)
	assert.Nil(t, err)
}

func TestTypeNameCorrectNonExported(t *testing.T) {
	file := ParseFileContent(`
		package test

		type fooBar struct{}
	`)

	checker := &TypeNameChecker{}
	err := checker.Check(file)
	assert.Nil(t, err)
}

func TestTypeNameIncorrect(t *testing.T) {
	file := ParseFileContent(`
		package test

		type foo_bar struct{}
	`)

	checker := &TypeNameChecker{}
	err := checker.Check(file)
	assert.NotNil(t, err)
}

func TestTypeNameNoType(t *testing.T) {
	file := ParseFileContent(`
		package test

		const foo_bar_1 = 21

		func foo_bar_2() {}
	`)

	checker := &TypeNameChecker{}
	err := checker.Check(file)
	assert.Nil(t, err)
}
