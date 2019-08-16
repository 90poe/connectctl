package manage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/90poe/connectctl/pkg/connect"
	"github.com/90poe/connectctl/pkg/manager"
	signals "github.com/90poe/connectctl/pkg/signal"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	definitionFiles []string
	syncPeriod      time.Duration
	allowPurge      bool
	directory       string
)

func manageConnectorsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "connectors",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long:  "",
		Run:   doManageConnectors,
	}

	cmd.Flags().StringArrayVarP(&definitionFiles, "files", "f", []string{}, "The connector definitions files (Required if --directory not specified)")
	_ = viper.BindPFlag("files", cmd.PersistentFlags().Lookup("files"))

	cmd.Flags().StringVarP(&directory, "directory", "d", "", "The directory containing the connector definitions files (Required if --files not specified)")
	_ = viper.BindPFlag("directory", cmd.PersistentFlags().Lookup("directory"))

	cmd.Flags().DurationVarP(&syncPeriod, "sync-period", "s", 5*time.Minute, "How often to sync with the connect cluster. Defaults to 5 minutes")
	_ = viper.BindPFlag("sync-period", cmd.PersistentFlags().Lookup("sync-period"))

	cmd.Flags().BoolVarP(&allowPurge, "allow-purge", "", false, "If true it will manage all connectors in a cluster. If connectors exist in the cluster that aren't specified in --files then the connectors will be deleted")
	_ = viper.BindPFlag("allow-purge", cmd.PersistentFlags().Lookup("allow-purge"))

	return cmd
}

func doManageConnectors(cmd *cobra.Command, args []string) {
	clusterLogger := log.WithField("cluster", clusterURL)
	clusterLogger.Debug("executing manage connectors command")

	err := checkConfig()
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error with configuration")
	}

	var source manager.ConnectorSource
	if definitionFiles != nil {
		source = filesSource(&definitionFiles)
	}
	if directory != "" {
		source = directorySource(&directory)
	}

	stopCh := signals.SetupSignalHandler()

	config := &manager.Config{
		ClusterURL: clusterURL,
		SyncPeriod: syncPeriod,
		Logger:     clusterLogger,
		AllowPurge: allowPurge,
	}
	clusterLogger.WithField("config", config).Trace("manage connectors confirguration")

	mngr, err := manager.NewConnectorsManager(config, source)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error creating connectors manager")
	}

	if err := mngr.Run(stopCh); err != nil {
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

func checkConfig() error {
	if len(definitionFiles) == 0 && directory == "" {
		return errors.New("you must supply a list of files using --files or a directory that contains files using --directory")
	}

	if len(definitionFiles) != 0 && directory != "" {
		return errors.New("you can't supply a list of files and a directory that contains files. Use --files OR --directory")
	}

	return nil
}
