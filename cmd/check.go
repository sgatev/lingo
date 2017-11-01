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
	// Matchers is a list of file matchers used to define
	// the files that will be checked.
	Matchers []struct {
		Type   string                 `yaml:"type"`
		Config map[string]interface{} `yaml:"config"`
	} `yaml:"matchers"`

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

		var matchers []file.Matcher
		for _, matcher := range config.Matchers {
			matchers = append(matchers,
				file.Get(matcher.Type, matcher.Config))
		}

		feeder := file.NewFeeder(matchers...)

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

		reports := map[string]*checker.Report{}

		for path := range files {
			reports[path] = &checker.Report{}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			file, err := parser.ParseFile(
				token.NewFileSet(),
				path,
				nil,
				parser.ParseComments)
			if err != nil {
				// TODO: handle error gracefully
				panic(err)
			}

			fc.Check(file, string(content), reports[path])
		}

		totalErrors := 0
		for path, report := range reports {
			if len(report.Errors) == 0 {
				continue
			}

			fmt.Println(path)
			for _, err := range report.Errors {
				fmt.Printf("\t- %s\n", err.Error())
			}
			fmt.Println()

			totalErrors += len(report.Errors)
		}
		fmt.Printf("%d violations found in %d files\n",
			totalErrors, len(reports))

		if totalErrors > 0 {
			os.Exit(1)
		}
	},
}

var configFile string
