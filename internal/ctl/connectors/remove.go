package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type removeConnectorsCmdParams struct {
	ClusterURL string
	Connectors []string
}

func removeConnectorCmd() *cobra.Command {
	params := &removeConnectorsCmdParams{}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove connectors from a connect cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doRemoveConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(removeCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(removeCmd, &params.Connectors)

	return removeCmd
}

func doRemoveConnectors(_ *cobra.Command, params *removeConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("removing connectors: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("remove connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}
	err = mngr.Remove(params.Connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error removing connectors")
	}

	clusterLogger.Info("removed connectors")
}
