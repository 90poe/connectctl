package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/90poe/connectctl/pkg/sources"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type addConnectorsCmdParams struct {
	ClusterURL string
	Files      []string
	Directory  string
	EnvVar     string
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
	ctl.AddDefinitionFilesFlags(addCmd, &params.Files, &params.Directory, &params.EnvVar)

	return addCmd
}

func doAddConnectors(cmd *cobra.Command, params *addConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("adding connectors")

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("add connectors configuration")

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

	connectors, err := source()
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error reading connector configuration from files")
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
