package connectors

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	signals "github.com/90poe/connectctl/pkg/signal"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type manageConnectorsCmdParams struct {
	ClusterURL  string
	Files       []string
	Directory   string
	SyncPeriod  time.Duration
	AllowPurge  bool
	AutoRestart bool
}

func manageConnectorsCmd() *cobra.Command {

	params := &manageConnectorsCmdParams{}

	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doManageConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(manageCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(manageCmd, &params.Files, &params.Directory)

	manageCmd.Flags().DurationVarP(&params.SyncPeriod, "sync-period", "s", 5*time.Minute, "How often to sync with the connect cluster. Defaults to 5 minutes")
	_ = viper.BindPFlag("sync-period", manageCmd.PersistentFlags().Lookup("sync-period"))

	manageCmd.Flags().BoolVarP(&params.AllowPurge, "allow-purge", "", false, "If true it will manage all connectors in a cluster. If connectors exist in the cluster that aren't specified in --files then the connectors will be deleted")
	_ = viper.BindPFlag("allow-purge", manageCmd.PersistentFlags().Lookup("allow-purge"))

	manageCmd.Flags().BoolVar(&params.AutoRestart, "auto-restart", false, "if supplied tasks that are failed with automatically be restarted")
	_ = viper.BindPFlag("auto-restart", manageCmd.PersistentFlags().Lookup("auto-restart"))

	return manageCmd
}

func doManageConnectors(_ *cobra.Command, params *manageConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Debug("executing manage connectors command")

	err := checkConfig(params)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error with configuration")
	}

	var source manager.ConnectorSource
	if params.Files != nil {
		source = filesSource(&params.Files)
	}
	if params.Directory != "" {
		source = directorySource(&params.Directory)
	}

	stopCh := signals.SetupSignalHandler()

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

	if err := mngr.Manage(source, stopCh); err != nil {
		clusterLogger.WithError(err).Fatalln("Error running connector manager")
	}

	clusterLogger.Info("finished executing manage connectors command")
}

func filesSource(files *[]string) manager.ConnectorSource {
	return func() ([]*connect.Connector, error) {
		connectors := make([]*connect.Connector, len(*files))

		for index, file := range *files {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, errors.Wrapf(err, "reading connector json %s", file)
			}

			connector := &connect.Connector{}
			err = json.Unmarshal(bytes, connector)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshalling connector from bytes")
			}

			connectors[index] = connector
		}

		return connectors, nil
	}
}

func directorySource(dir *string) manager.ConnectorSource {
	return func() ([]*connect.Connector, error) {
		var files []string

		err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".json" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "list connector files in directory %s", *dir)
		}

		return filesSource(&files)()
	}
}

func checkConfig(params *manageConnectorsCmdParams) error {
	if len(params.Files) == 0 && params.Directory == "" {
		return errors.New("you must supply a list of files using --files or a directory that contains files using --directory")
	}

	if len(params.Files) != 0 && params.Directory != "" {
		return errors.New("you can't supply a list of files and a directory that contains files. Use --files OR --directory")
	}

	return nil
}
