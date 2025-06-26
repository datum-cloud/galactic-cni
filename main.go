package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/datum-cloud/galactic/cli"
	"github.com/datum-cloud/galactic/cni"
)

func main() {
	basename := filepath.Base(os.Args[0])
	args := os.Args[1:]

	var cmd *cobra.Command

	switch basename {
	case "galactic-cli":
		cmd = cli.NewCommand()
	case "galactic-cni":
		cmd = cni.NewCommand()
	default:
		log.Fatalf("Unknown binary name: %s. Should be one of galactic-cli or galactic-cni.", basename)
	}

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Execution failed: %v", err)
	}
}
