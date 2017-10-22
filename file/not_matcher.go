package file

type notMatcher struct {
	matcher Matcher
}

// NotMatcher creates a new Matcher that reverses the decision
// of `matcher`.
func NotMatcher(matcher Matcher) Matcher {
	return &notMatcher{
		matcher: matcher,
	}
}

// Matches implements the Matcher interface.
func (m *notMatcher) Matches(path string) bool {
	return !m.matcher.Matches(path)
}
