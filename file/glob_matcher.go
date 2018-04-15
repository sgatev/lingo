package file

import (
	"github.com/mattn/go-zglob"
	"github.com/uber-go/mapdecode"
)

func init() {
	must(Register("glob", GlobMatcher))
}

// GlobMatcherConfig describes the configuration of a GlobMatcher.
type GlobMatcherConfig struct {

	// Pattern is the glob pattern used by the matcher.
	Pattern string `yaml:"pattern"`
}

type globMatcher struct {
	pattern string
}

// GlobMatcher creates a new Matcher that accepts files based on
// glob pattern.
func GlobMatcher(configData interface{}) Matcher {
	var config GlobMatcherConfig
	if err := mapdecode.Decode(&config, configData); err != nil {
		return nil
	}

	return &globMatcher{
		pattern: config.Pattern,
	}
}

// Matches implements the Matcher interface.
func (m *globMatcher) Matches(path string) bool {
	ok, err := zglob.Match(m.pattern, path)
	return err == nil && ok
}
