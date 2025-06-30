package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datum-cloud/galactic/cli"
	"github.com/datum-cloud/galactic/cni"
	"github.com/datum-cloud/galactic/debug"
)

func main() {
	basename := filepath.Base(os.Args[0])
	args := os.Args[1:]

	var cmd *cobra.Command

	switch {
	case strings.HasPrefix(basename, "galactic-cli"):
		cmd = cli.NewCommand()
	case strings.HasPrefix(basename, "galactic-cni"):
		cmd = cni.NewCommand()
	default:
		log.Fatalf("Unknown binary name: %s. Should be one of galactic-cli or galactic-cni.", basename)
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version details",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(debug.Version())
		},
	}
	cmd.AddCommand(versionCmd)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Execution failed: %v", err)
	}
}
