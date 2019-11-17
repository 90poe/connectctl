package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/90poe/connectctl/pkg/sources"
	"github.com/pkg/errors"

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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doAddConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(addCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(addCmd, &params.Files, &params.Directory, &params.EnvVar)

	return addCmd
}

func doAddConnectors(cmd *cobra.Command, params *addConnectorsCmdParams) error {
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
		return errors.New("error finding connector definitions from parameters")
	}

	connectors, err := source()
	if err != nil {
		return errors.Wrap(err, "error reading connector configuration from files")
	}

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}
	err = mngr.Add(connectors)
	if err != nil {
		return errors.Wrap(err, "error creating connectors")
	}
	clusterLogger.Infof("added connectors")
	return nil
}
