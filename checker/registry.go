package checker

import "fmt"

// NodeCheckerConstructor constructs NodeChecker instances.
type NodeCheckerConstructor func() NodeChecker

// Register adds `checker` to the registry.
func Register(constructor NodeCheckerConstructor) error {
	slug := constructor().Slug()

	if _, ok := registry[slug]; ok {
		return fmt.Errorf("checker already registered: " + slug)
	}

	registry[slug] = constructor

	return nil
}

// Get returns the NodeChecker referenced by a `slug`.
func Get(slug string) NodeChecker {
	constructor, ok := registry[slug]
	if !ok {
		return nil
	}

	return constructor()
}

var registry = map[string]NodeCheckerConstructor{}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
