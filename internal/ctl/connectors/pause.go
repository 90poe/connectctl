package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

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
		Run: func(cmd *cobra.Command, _ []string) {
			doPauseConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(pauseCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(pauseCmd, &params.Connectors)

	return pauseCmd
}

func doPauseConnectors(_ *cobra.Command, params *pauseConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("pausing connectors: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("pause connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	err = mngr.Pause(params.Connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatal("error pausing connectors")
	}

	clusterLogger.Info("connectors paused successfully")
}
