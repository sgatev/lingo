package checker_test

import (
	"go/ast"
	"go/parser"

	"golang.org/x/tools/go/loader"
)

// ParseFileContent parses `content` and returns an AST.
func ParseFileContent(content string) *ast.File {
	config := &loader.Config{
		ParserMode: parser.ParseComments,
	}

	file, err := config.ParseFile("test.go", content)
	if err != nil {
		panic("could not parse file content")
	}

	return file
}
