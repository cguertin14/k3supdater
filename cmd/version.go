package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitCommit
	// The latest git commit hash
	GitCommit string

	// BuildDate
	// The date when the build was made
	BuildDate string

	// Version
	// The version of the app
	Version string
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show k3supdater's installed version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(
				"Version: %s, GitCommit: %s, BuildDate: %s\n",
				Version, GitCommit, BuildDate,
			)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
