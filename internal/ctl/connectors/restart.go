package connectors

import (
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type restartConnectorsCmdParams struct {
	ClusterURL string
	Connectors []string
}

// Command creates the tehe management commands
func restartConnectorsCmd() *cobra.Command {
	params := &restartConnectorsCmdParams{}

	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart connectors in a cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doRestartConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(restartCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(restartCmd, &params.Connectors)

	return restartCmd
}

func doRestartConnectors(_ *cobra.Command, params *restartConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("restarting connectors: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("restart connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	err = mngr.Restart(params.Connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatal("error restarting connectors")
	}

	clusterLogger.Info("connectors restarted successfully")
}
