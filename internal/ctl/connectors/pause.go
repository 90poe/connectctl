package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pauseConnectorsCmdParams struct {
	ClusterURL string
	Connectors []string
}

func pauseConnectorsCmd() *cobra.Command {
	params := &pauseConnectorsCmdParams{}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause connectors in a cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doPauseConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(pauseCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(pauseCmd, &params.Connectors)

	return pauseCmd
}

func doPauseConnectors(_ *cobra.Command, params *pauseConnectorsCmdParams) error {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("pausing connectors: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("pause connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}

	err = mngr.Pause(params.Connectors)
	if err != nil {
		return errors.Wrap(err, "error pausing connectors")
	}

	clusterLogger.Info("connectors paused successfully")
	return nil
}
