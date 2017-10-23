package cmd

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/s2gatev/lingo/checker"
	"github.com/s2gatev/lingo/file"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(Check)
}

// Check is a command handler that checks the lingo in a directory
// for violations.
var Check = &cobra.Command{
	Use:   "check",
	Short: "Check the lingo of all files in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		feeder := file.NewFeeder(
			file.GlobMatcher("**/*.go"),
			file.NotMatcher(file.GlobMatcher("**/vendor/**/*")),
			file.NotMatcher(file.GlobMatcher("**/*_test.go")))

		// TODO: parse slugs from config file
		slugs := []string{
			"local_return",
			"multi_word_ident_name",
			"exported_ident_doc",
		}

		fc := checker.NewFileChecker()
		for _, slug := range slugs {
			c := checker.Get(slug)
			if c == nil {
				// TODO: handle error gracefully
				panic("unknown checker: " + slug)
			}

			fc.Register(c)
		}

		files, err := feeder.Feed(args[0])
		if err != nil {
			// TODO: handle error gracefully
			panic(err)
		}

		var report checker.Report

		fset := token.NewFileSet()
		for path := range files {
			file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				// TODO: handle error gracefully
				panic(err)
			}

			fc.Check(file, &report)
		}

		for _, err := range report.Errors {
			fmt.Println(err.Error())
		}

		if len(report.Errors) > 0 {
			os.Exit(1)
		}
	},
}
