package checker

import "fmt"

// Register adds `checker` to the registry.
func Register(checker NodeChecker) error {
	slug := checker.Slug()

	if _, ok := registry[slug]; ok {
		return fmt.Errorf("checker already registered: " + slug)
	}

	registry[slug] = checker

	return nil
}

// Get returns the NodeChecker referenced by a `slug`.
func Get(slug string) NodeChecker {
	return registry[slug]
}

var registry = map[string]NodeChecker{}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
