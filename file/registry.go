package file

import "fmt"

// MatcherConstructor constructs Matcher instances.
type MatcherConstructor func(configData interface{}) Matcher

// Register adds a matcher to the registry.
func Register(slug string, constructor MatcherConstructor) error {
	if _, ok := registry[slug]; ok {
		return fmt.Errorf("matcher already registered: " + slug)
	}

	registry[slug] = constructor

	return nil
}

// Get returns the Matcher referenced by a `slug`.
func Get(slug string, config interface{}) Matcher {
	constructor, ok := registry[slug]
	if !ok {
		return nil
	}

	return constructor(config)
}

var registry = map[string]MatcherConstructor{}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
