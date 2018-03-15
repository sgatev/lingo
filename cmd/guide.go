package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/template"
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
		dir, err := ioutil.TempDir("", "lingo")
		if err != nil {
			return
		}

		guide, err := os.Create(filepath.Join(dir, "guide.html"))
		if err != nil {
			return
		}
		defer guide.Close()

		if err := guideTemplate.Execute(guide, nil); err != nil {
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

var guideTemplate = template.Must(template.New("html").Parse(guideContent))

const guideContent = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head>
	<body>
	</body>
</html>
`
