package manage

import (
	"encoding/json"
	"io/ioutil"
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
)

func manageConnectorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectors",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long:  "",
		Run:   doManageConnectors,
	}

	cmd.Flags().StringArrayVarP(&definitionFiles, "files", "f", []string{}, "The connector definitions files (Required)")
	cmd.MarkFlagRequired("files")
	viper.BindPFlag("files", cmd.PersistentFlags().Lookup("files"))

	cmd.Flags().DurationVarP(&syncPeriod, "sync-period", "s", 5*time.Minute, "How often to sync with the connect cluster. Defaults to 5 minutes")
	viper.BindPFlag("sync-period", cmd.PersistentFlags().Lookup("sync-period"))

	return cmd
}

func doManageConnectors(cmd *cobra.Command, args []string) {
	clusterLogger := log.WithField("cluster", clusterURL)
	clusterLogger.Debug("executing manage connectors command")

	connectors, err := loadDefinitions(definitionFiles)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error reading connectors files")
	}

	stopCh := signals.SetupSignalHandler()

	mngr, err := manager.NewConnectorsManager(clusterURL, connectors, syncPeriod, clusterLogger)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("Error creating connectors manager")
	}

	if err := mngr.Run(stopCh); err != nil {
		clusterLogger.WithError(err).Fatalln("Error running connector manager")
	}

	clusterLogger.Info("finished executing manage connectors command")
}

func loadDefinitions(files []string) ([]*connect.Connector, error) {
	connectors := make([]*connect.Connector, len(files))

	for index, file := range files {
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
