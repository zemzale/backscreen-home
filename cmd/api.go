package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API",
	Long:  `Start the API to get stored currency exchange rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement the API
		return errors.New("not implemented")
	},
}
