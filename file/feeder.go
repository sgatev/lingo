package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

		if strings.HasSuffix(dir, recPathSuffix) {
			f.feedDirRecursive(strings.TrimSuffix(dir, recPathSuffix), files)
		} else {
			f.feedDir(dir, files)
		}
	}()

	return files, nil
}

const recPathSuffix = "/..."

func (f *Feeder) feedDir(dir string, paths chan<- string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())
		if f.matches(path) {
			paths <- path
		}
	}
}

func (f *Feeder) feedDirRecursive(dir string, paths chan<- string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if f.matches(path) {
			paths <- path
		}

		return nil
	})
}

func (f *Feeder) matches(path string) bool {
	for _, matcher := range f.matchers {
		if !matcher.Matches(path) {
			return false
		}
	}

	return true
}
