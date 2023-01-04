package cmd

import (
	"fmt"

	dbutils "github.com/infratographer/lmi/internal/storage/sql/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"

	"github.com/infratographer/lmi/internal/storage/sql/migrations"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Executes a database migration",
	Long:  `Executes a database migration based on the current version of the database.`,
	RunE:  migrate,
}

//nolint:gochecknoinits // This is a Cobra generated file
func init() {
	rootCmd.AddCommand(migrateCmd)

	v := viper.GetViper()
	flags := migrateCmd.Flags()

	crdbx.MustViperFlags(v, flags)
}

func migrate(cmd *cobra.Command, args []string) error {
	v := viper.GetViper()

	// Initialize logger
	l := initLogger()

	// Initialize database connection
	dbconn, err := dbutils.GetDBConnection(v, defaultDBName, false)
	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}

	l.Info("executing migrations")

	return migrations.Migrate(dbconn)
}
