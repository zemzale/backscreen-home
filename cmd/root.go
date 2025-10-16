package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zemzale/backscreen-home/adapter/database"
	"github.com/zemzale/backscreen-home/storage"
)

var store *storage.Client

var rootCmd = &cobra.Command{
	Use:   "backscreen-home",
	Short: "Currency exchange rates service",
	Long:  `Currency exchange rates service that has two main modes. Fetch rates and server an API to get those.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		logger := slog.With("component", "root")

		logger.DebugContext(ctx, "Connecting to database")
		db, err := database.New(&database.Config{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Database: viper.GetString("database.database"),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		logger.DebugContext(ctx, "Creating storage client")
		store = storage.New(db)

		if err := store.Migrate(ctx); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Setup Viper to work with environment variables
	viper.SetEnvPrefix("BACKSCREEN")
	viper.AutomaticEnv()

	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(syncCmd)
}
