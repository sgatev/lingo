package checker

// Register adds `checker` to the registry.
func Register(checker NodeChecker) {
	registry[checker.Slug()] = checker
}

// Get returns the NodeChecker referenced by a `slug`.
func Get(slug string) NodeChecker {
	return registry[slug]
}

var registry = map[string]NodeChecker{}
