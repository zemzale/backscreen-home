package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "backscreen-home",
	Short: "Currency exchange rates service",
	Long:  `Currency exchange rates service that has two main modes. Fetch rates and server an API to get those.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(syncCmd)
}
