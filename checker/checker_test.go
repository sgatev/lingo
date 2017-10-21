package checker_test

import (
	"go/ast"

	"golang.org/x/tools/go/loader"
)

func ParseFileContent(content string) *ast.File {
	config := &loader.Config{}

	file, err := config.ParseFile("test.go", content)
	if err != nil {
		panic("could not parse file content")
	}

	return file
}
