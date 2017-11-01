package file

import "github.com/uber-go/mapdecode"

func init() {
	must(Register("not", NotMatcher))
}

// NotMatcherConfig describes the configuration of a NotMatcher.
type NotMatcherConfig struct {

	// Type is the type of the negated matcher.
	Type string `yaml:"type"`

	// Config is the configuration of the negated matcher.
	Config interface{} `yaml:"config"`
}

type notMatcher struct {
	matcher Matcher
}

// NotMatcher creates a new Matcher that reverses the decision
// of matcher.
func NotMatcher(configData interface{}) Matcher {
	var config NotMatcherConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &notMatcher{
		matcher: Get(config.Type, config.Config),
	}
}

// Matches implements the Matcher interface.
func (m *notMatcher) Matches(path string) bool {
	return !m.matcher.Matches(path)
}
