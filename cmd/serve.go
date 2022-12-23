/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/google/uuid"
	apiv1 "github.com/infratographer/fertilesoil/api/v1"
	appv1 "github.com/infratographer/fertilesoil/app/v1"
	appv1sql "github.com/infratographer/fertilesoil/app/v1/sql"
	clientv1 "github.com/infratographer/fertilesoil/client/v1"
	cv1nats "github.com/infratographer/fertilesoil/client/v1/nats"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/ginx"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/viperx"
	"go.uber.org/zap"

	"github.com/infratographer/lmi/internal/reconciler"
	dbutils "github.com/infratographer/lmi/internal/storage/sql/utils"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: serve,
}

const (
	defaultListenAddr = ":8080"
	defaultDBName     = "permissions"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	v := viper.GetViper()
	flags := serveCmd.Flags()

	crdbx.MustViperFlags(v, flags)
	ginx.MustViperFlags(v, flags, defaultListenAddr)
	loggingx.MustViperFlags(v, flags)

	flags.String("base-directory-id", "", "ID of the base directory for this lmi instance")
	viperx.MustBindFlag(v, "base_directory_id", flags.Lookup("base-directory-id"))

	flags.String("nats-url", "", "NATS URL")
	viperx.MustBindFlag(v, "nats.url", flags.Lookup("nats-url"))

	flags.String("nats-directories-subjects", "infratographer.events.directories",
		"NATS subject to register to directory events")
	viperx.MustBindFlag(v, "nats.directories_subjects", flags.Lookup("nats-directories-subjects"))

	flags.String("nats-nkey", "", "path to nkey file")
	viperx.MustBindFlag(v, "nats.nkey", flags.Lookup("nats-nkey"))
}

func initLogger() *zap.Logger {
	sl := loggingx.InitLogger("lmi", loggingx.Config{
		Debug:  viper.GetBool("debug"),
		Pretty: viper.GetBool("pretty"),
	})

	return sl.Desugar()
}

func serve(cmd *cobra.Command, args []string) error {
	v := viper.GetViper()

	// Initialize logger
	logger := initLogger()

	// Initialize database connection
	dbconn, err := dbutils.GetDBConnection(v, defaultDBName, false)
	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}

	// Initialize app storage
	appStore := appv1sql.New(dbconn)

	// Initialize NATS connection
	opts := []nats.Option{
		nats.Name("lmi"),
	}

	opt, err := nats.NkeyOptionFromSeed(viper.GetString("nats.nkey"))
	if err != nil {
		return fmt.Errorf("failed to load nkey: %w", err)
	}

	opts = append(opts, opt)

	natsconn, err := nats.Connect(v.GetString("nats.url"), opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}

	// Create NATS directory subscriber
	watcher, err := cv1nats.NewSubscriber(natsconn, viper.GetString("nats.directories_subjects"))
	if err != nil {
		return fmt.Errorf("failed to create nats subscriber: %w", err)
	}

	// Create dirclient
	// TODO(jaosorior): Make this configurable.
	// We'd need parameters for the dirclient to be able to connect to the server.
	dirclient := clientv1.NewHTTPClient(nil)

	// Initialize our reconciler
	r := reconciler.NewReconciler()

	// Get base directory
	rawID := v.GetString("base_directory_id")

	baseDirID, err := uuid.Parse(rawID)
	if err != nil {
		return fmt.Errorf("failed to parse base directory id: %w", err)
	}

	ctrl, err := appv1.NewController(
		apiv1.DirectoryID(baseDirID),
		appv1.WithStorage(appStore),
		appv1.WithWatcher(watcher),
		appv1.WithClient(dirclient),
		appv1.WithReconciler(r),
	)
	if err != nil {
		return fmt.Errorf("failed to create directory controller: %w", err)
	}

	ctx := cmd.Context()

	go func() {
		if err := ctrl.Run(ctx); err != nil {
			logger.Fatal("failed to run controller", zap.Error(err))
		}
	}()

	// Run permissions API server

	return nil
}
