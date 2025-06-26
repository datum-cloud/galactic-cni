package cni

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "galactic-cni",
		Short: "Galactic CNI plugin",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Running CNI logic with args: %v\n", args)
		},
	}
}
