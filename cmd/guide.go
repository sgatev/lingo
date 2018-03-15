package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"

	"github.com/alecthomas/template"
	"github.com/s2gatev/lingo/checker"
	"github.com/spf13/cobra"
)

func init() {
	Guide.PersistentFlags().StringVar(
		&configFile, "config", defaultConfigFilename, "config file")

	Root.AddCommand(Guide)
}

// Guide is a command handler that displays a guidebook of rules applicable
// for the current project.
var Guide = &cobra.Command{
	Use:   "guide",
	Short: "Read a guide with the lingo of the project",
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

		var checkers []checker.NodeChecker
		for slug, config := range config.Checkers {
			c := checker.Get(slug, config)
			if c == nil {
				// TODO: handle error gracefully
				panic("unknown checker: " + slug)
			}

			checkers = append(checkers, c)
		}

		var items []guideItem
		for _, checker := range checkers {
			items = append(items, guideItem{
				Title:       checker.Title(),
				Description: checker.Description(),
			})
		}

		dir, err := ioutil.TempDir("", "lingo")
		if err != nil {
			// TODO: handle error gracefully
			panic(err)
		}

		guide, err := os.Create(filepath.Join(dir, "guide.html"))
		if err != nil {
			// TODO: handle error gracefully
			panic(err)
		}
		defer guide.Close()

		if err := guideTemplate.Execute(guide, items); err != nil {
			return
		}

		if err := openBrowser("file://" + guide.Name()); err != nil {
			return
		}
	},
}

// openBrowser tries to open the URL in a browser.
func openBrowser(url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}

	return exec.Command(args[0], append(args[1:], url)...).Run()
}

type guideItem struct {

	// Title is the title of the item.
	Title string

	// Description is the detailed description of the item.
	Description string
}

var guideTemplate = template.Must(template.New("html").Parse(guideContent))

const guideContent = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head>
	<body>
		{{range .}}
			<div>
				<h2>{{.Title}}</h2>
				<p>{{.Description}}</p>
			</div>
		{{end}}
	</body>
</html>
`
