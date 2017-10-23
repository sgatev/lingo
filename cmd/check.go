package cmd

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/s2gatev/lingo/checker"
	"github.com/s2gatev/lingo/file"
	"github.com/spf13/cobra"
)

func init() {
	Check.PersistentFlags().StringVar(&configFile, "config", "", "config file")

	Root.AddCommand(Check)
}

// Config describes the lingo check config file structure.
type Config struct {

	// Checkers is a map[checker_slug]checker_config of checkers
	// that need to be executed.
	Checkers map[string]interface{} `yaml:"checkers"`
}

// Check is a command handler that checks the lingo in a directory
// for violations.
var Check = &cobra.Command{
	Use:   "check",
	Short: "Check the lingo of all files in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			// TODO: handle error gracefully
			panic(err)
		}

		var config Config
		if err := yaml.Unmarshal(configData, &config); err != nil {
			// TODO: handle error gracefully
			panic(err)
		}

		feeder := file.NewFeeder(
			file.GlobMatcher("**/*.go"),
			file.NotMatcher(file.GlobMatcher("**/vendor/**/*")),
			file.NotMatcher(file.GlobMatcher("**/*_test.go")))

		fc := checker.NewFileChecker()
		for slug := range config.Checkers {
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

var configFile string
