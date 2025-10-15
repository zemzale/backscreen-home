package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync currency exchange rates",
	Long:  `Sync currency exchange rates from the source to the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement the sync
		return errors.New("not implemented")
	},
}
