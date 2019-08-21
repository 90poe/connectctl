package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type addConnectorsCmdParams struct {
	ClusterURL string
	Files      []string
	Directory  string
}

func addConnectorCmd() *cobra.Command {

	params := &addConnectorsCmdParams{}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add connectors to a connect cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doAddConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(addCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(addCmd, &params.Files, &params.Directory)

	return addCmd
}

func doAddConnectors(_ *cobra.Command, params *addConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("adding connectors")

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("add connectors configuration")

	var source manager.ConnectorSource
	if params.Files != nil {
		source = filesSource(&params.Files)
	}
	if params.Directory != "" {
		source = directorySource(&params.Directory)
	}

	connectors, err := source()
	if err != nil {
		clusterLogger.WithError(err).Error("error reading connector configuration from files")
	}

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}
	err = mngr.Add(connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors")
	}
	clusterLogger.Infof("added connectors")
}
