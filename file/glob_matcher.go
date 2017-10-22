package file

import zglob "github.com/mattn/go-zglob"

type globMatcher struct {
	pattern string
}

// GlobMatcher creates a new Matcher that accepts files based on
// glob `pattern`.
func GlobMatcher(pattern string) Matcher {
	return &globMatcher{
		pattern: pattern,
	}
}

// Matches implements the Matcher interface.
func (m *globMatcher) Matches(path string) bool {
	ok, err := zglob.Match(m.pattern, path)
	return err == nil && ok
}
