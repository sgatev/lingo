package cli

import (
	"fmt"
	"os"
)

// ExitOK prints a formatted message and exits successfully.
func ExitOK(message string, args ...interface{}) {
	print(message, args...)
	os.Exit(0)
}

// ExitError prints a formatted message and exits with error.
func ExitError(message string, args ...interface{}) {
	print(message, args...)
	os.Exit(1)
}

func print(message string, args ...interface{}) {
	fmt.Printf("lingo: %s\n", fmt.Sprintf(message, args...))
}
