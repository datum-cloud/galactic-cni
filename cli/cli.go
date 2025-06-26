package cli

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "galactic-cli",
		Short: "Galactic CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Running CLI logic with args: %v\n", args)
		},
	}
}
