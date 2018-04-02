package checker_test

import (
	"go/ast"
	"go/parser"
	"go/token"

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

// ParseFileContentInSet parses `content` adding it to the file set and returns an AST.
func ParseFileContentInSet(fileSet *token.FileSet, content string) *ast.File {
	config := &loader.Config{
		Fset:       fileSet,
		ParserMode: parser.ParseComments,
	}

	file, err := config.ParseFile("test.go", content)
	if err != nil {
		panic("could not parse file content")
	}

	return file
}
