package cmd

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
	Checkers map[string]map[string]interface{} `yaml:"checkers"`
}

const defaultConfigFilename = "lingo.yml"

var configFile string
