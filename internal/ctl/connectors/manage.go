package connectors

import (
	"context"
	"time"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/healthcheck"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"
	signals "github.com/90poe/connectctl/pkg/signal"
	"github.com/90poe/connectctl/pkg/sources"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type manageConnectorsCmdParams struct {
	ClusterURL         string
	Files              []string
	Directory          string
	EnvVar             string
	SyncPeriod         time.Duration
	AllowPurge         bool
	AutoRestart        bool
	RunOnce            bool
	EnableHealthCheck  bool
	HealthCheckAddress string
}

func manageConnectorsCmd() *cobra.Command {
	params := &manageConnectorsCmdParams{
		SyncPeriod:         5 * time.Minute,
		HealthCheckAddress: ":9000",
	}

	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long: `This command will add/delete/update connectors in a destination
Kafa Connect cluster based on a list of desired connectors which are specified
as a list of files or all files in a directory. The command runs continuously and
will sync desired state with actual state based on the --sync-period flag. But
if you specify --once then it will sync once and then exit.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doManageConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(manageCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(manageCmd, &params.Files, &params.Directory, &params.EnvVar)

	manageCmd.Flags().DurationVarP(&params.SyncPeriod, "sync-period", "s", params.SyncPeriod, "How often to sync with the connect cluster")
	_ = viper.BindPFlag("sync-period", manageCmd.PersistentFlags().Lookup("sync-period"))

	manageCmd.Flags().BoolVarP(&params.AllowPurge, "allow-purge", "", false, "If true it will manage all connectors in a cluster. If connectors exist in the cluster that aren't specified in --files then the connectors will be deleted")
	_ = viper.BindPFlag("allow-purge", manageCmd.PersistentFlags().Lookup("allow-purge"))

	manageCmd.Flags().BoolVar(&params.AutoRestart, "auto-restart", false, "if supplied tasks that are failed with automatically be restarted")
	_ = viper.BindPFlag("auto-restart", manageCmd.PersistentFlags().Lookup("auto-restart"))

	manageCmd.Flags().BoolVar(&params.RunOnce, "once", false, "if supplied sync will run once and command will exit")
	_ = viper.BindPFlag("once", manageCmd.PersistentFlags().Lookup("once"))

	manageCmd.Flags().BoolVar(&params.EnableHealthCheck, "healthcheck-enable", false, "if supplied a healthcheck via http will be enabled")
	_ = viper.BindPFlag("healthcheck-enable", manageCmd.PersistentFlags().Lookup("healthcheck-enable"))

	manageCmd.Flags().StringVar(&params.HealthCheckAddress, "healthcheck-address", params.HealthCheckAddress, "if enabled the healthchecks ('/live' and '/ready') will be available from this address")
	_ = viper.BindPFlag("healthcheck-address", manageCmd.PersistentFlags().Lookup("healthcheck-address"))

	return manageCmd
}

func doManageConnectors(cmd *cobra.Command, params *manageConnectorsCmdParams) error {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Debug("executing manage connectors command")

	err := checkConfig(params)
	if err != nil {
		return errors.Wrap(err, "Error with configuration")
	}

	source, err := findSource(params.Files, params.Directory, params.EnvVar, cmd)

	if err != nil {
		return err
	}

	config := &manager.Config{
		ClusterURL:  params.ClusterURL,
		SyncPeriod:  params.SyncPeriod,
		AllowPurge:  params.AllowPurge,
		AutoRestart: params.AutoRestart,
		Version:     version.Version,
	}
	clusterLogger.WithField("config", config).Trace("manage connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		return errors.Wrap(err, "Error creating connectors manager")
	}

	ctx := context.Background()

	if params.EnableHealthCheck {
		healthCheckHandler := healthcheck.New(mngr)

		go func() {
			err := healthCheckHandler.Start(params.HealthCheckAddress)
			if err != nil {
				clusterLogger.WithError(err, "Error starting healthcheck")
			}
		}()

		// nolint
		defer healthCheckHandler.Shutdown(ctx)
	}

	if params.RunOnce {
		if err := mngr.Sync(source); err != nil {
			return errors.Wrap(err, "Error running connector sync")
		}
	} else {
		stopCh := signals.SetupSignalHandler()

		if err := mngr.Manage(source, stopCh); err != nil {
			return errors.Wrap(err, "Error running connector manager")
		}
	}

	clusterLogger.Info("finished executing manage connectors command")
	return nil
}

func findSource(files []string, directory, envar string, cmd *cobra.Command) (manager.ConnectorSource, error) {
	switch {
	case files != nil && len(files) > 0:
		if len(files) == 1 && files[0] == "-" {
			return sources.StdIn(cmd.InOrStdin()), nil
		}
		return sources.Files(files), nil

	case directory != "":
		return sources.Directory(directory), nil
	case envar != "":
		return sources.EnvVarValue(envar), nil
	}
	return nil, errors.New("error finding connector definitions from parameters")
}

func checkConfig(params *manageConnectorsCmdParams) error {
	paramsSet := 0

	if len(params.Files) != 0 {
		paramsSet++
	}
	if params.Directory != "" {
		paramsSet++
	}
	if params.EnvVar != "" {
		paramsSet++
	}

	if paramsSet == 1 {
		return nil
	}

	return errors.New("you must supply a list of files using --files or a directory that contains files using --directory or an environmental whose value is a JSON serialised connector or array of connectors")
}
