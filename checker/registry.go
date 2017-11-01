package checker

import "fmt"

// NodeCheckerConstructor constructs NodeChecker instances.
type NodeCheckerConstructor func(configData interface{}) NodeChecker

// Register adds `checker` to the registry.
func Register(slug string, constructor NodeCheckerConstructor) error {
	if _, ok := registry[slug]; ok {
		return fmt.Errorf("checker already registered: " + slug)
	}

	registry[slug] = constructor

	return nil
}

// Get returns the NodeChecker referenced by a `slug`.
func Get(slug string, config interface{}) NodeChecker {
	constructor, ok := registry[slug]
	if !ok {
		return nil
	}

	return constructor(config)
}

var registry = map[string]NodeCheckerConstructor{}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
