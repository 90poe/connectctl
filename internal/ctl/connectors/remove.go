package connectors

import (
	"fmt"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doRemoveConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(removeCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(removeCmd, &params.Connectors)

	return removeCmd
}

func doRemoveConnectors(_ *cobra.Command, params *removeConnectorsCmdParams) error {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("removing connectors: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("remove connectors configuration")

	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(params.ClusterURL, userAgent)
	if err != nil {
		return errors.Wrap(err, "error creating connect client")
	}

	mngr, err := manager.NewConnectorsManager(client, config)
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}
	err = mngr.Remove(params.Connectors)
	if err != nil {
		return errors.Wrap(err, "error removing connectors")
	}

	clusterLogger.Info("removed connectors")
	return nil
}
