package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "k3supdater",
		Short: "k3s ansible playbook version updater tool",
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
