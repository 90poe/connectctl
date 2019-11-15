package connectors

import (
	"time"

	"github.com/90poe/connectctl/internal/ctl"
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
	ClusterURL  string
	Files       []string
	Directory   string
	EnvVar      string
	SyncPeriod  time.Duration
	AllowPurge  bool
	AutoRestart bool
	RunOnce     bool
}

func manageConnectorsCmd() *cobra.Command {
	params := &manageConnectorsCmdParams{}

	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long: `This command will add/delete/update connectors in a destination
Kafa Connect cluster based on a list of desired connectors which are specified
as a list of files or all files in a directory. The command runs continuously and
will sync desired state with actual state based on the --sync-period flag. But
if you specify --once then it will sync once and then exit.`,
		Run: func(cmd *cobra.Command, _ []string) {
			doManageConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(manageCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(manageCmd, &params.Files, &params.Directory, &params.EnvVar)

	manageCmd.Flags().DurationVarP(&params.SyncPeriod, "sync-period", "s", 5*time.Minute, "How often to sync with the connect cluster. Defaults to 5 minutes")
	_ = viper.BindPFlag("sync-period", manageCmd.PersistentFlags().Lookup("sync-period"))

	manageCmd.Flags().BoolVarP(&params.AllowPurge, "allow-purge", "", false, "If true it will manage all connectors in a cluster. If connectors exist in the cluster that aren't specified in --files then the connectors will be deleted")
	_ = viper.BindPFlag("allow-purge", manageCmd.PersistentFlags().Lookup("allow-purge"))

	manageCmd.Flags().BoolVar(&params.AutoRestart, "auto-restart", false, "if supplied tasks that are failed with automatically be restarted")
	_ = viper.BindPFlag("auto-restart", manageCmd.PersistentFlags().Lookup("auto-restart"))

	manageCmd.Flags().BoolVar(&params.RunOnce, "once", false, "if supplied sync will run once and command will exit")
	_ = viper.BindPFlag("once", manageCmd.PersistentFlags().Lookup("once"))

	return manageCmd
}

func doManageConnectors(cmd *cobra.Command, params *manageConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Debug("executing manage connectors command")

	err := checkConfig(params)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error with configuration")
	}

	var source manager.ConnectorSource

	switch {
	case params.Files != nil:
		if len(params.Files) == 1 && params.Files[0] == "-" {
			source = sources.StdIn(cmd.InOrStdin())
		} else {
			source = sources.Files(params.Files)
		}
	case params.Directory != "":
		source = sources.Directory(params.Directory)
	case params.EnvVar != "":
		source = sources.EnvVarValue(params.EnvVar)
	default:
		clusterLogger.Fatalln("error finding connector definitions from parameters")
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
		clusterLogger.WithError(err).Fatalln("Error creating connectors manager")
	}

	if params.RunOnce {
		if err := mngr.Sync(source); err != nil {
			clusterLogger.WithError(err).Fatalln("Error running connector sync")
		}
	} else {
		stopCh := signals.SetupSignalHandler()

		if err := mngr.Manage(source, stopCh); err != nil {
			clusterLogger.WithError(err).Fatalln("Error running connector manager")
		}
	}

	clusterLogger.Info("finished executing manage connectors command")
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
