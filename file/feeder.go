package file

import (
	"os"
	"path/filepath"
)

// Matcher matches a file based on a path criteria.
type Matcher interface {

	// Matches matches a file based on path.
	Matches(path string) bool
}

// Feeder feeds files.
type Feeder struct {
	matchers []Matcher
}

// NewFeeder creates a new Feeder that feeds files accepted by
// all `matchers`.
func NewFeeder(matchers ...Matcher) *Feeder {
	return &Feeder{
		matchers: matchers,
	}
}

// Feed feeds all files starting from a `root` directory to a
// chan of paths.
// If the error return value is not nil then the chan return
// value is nil.
func (f *Feeder) Feed(root string) (<-chan string, error) {
	dir, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	files := make(chan string)

	go func() {
		defer close(files)

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			for _, matcher := range f.matchers {
				if !matcher.Matches(path) {
					return nil
				}
			}

			files <- path

			return nil
		})
	}()

	return files, nil
}
