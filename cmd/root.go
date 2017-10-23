package cmd

import "github.com/spf13/cobra"

// Root is a dummy command handler.
var Root = &cobra.Command{
	Use:   "lingo",
	Short: "Lingo helps you define and enforce project-specific Go lingo",
}
