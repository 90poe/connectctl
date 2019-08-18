package connectors

import (
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	connectors []string
)

// Command creates the tehe management commands
func restartConnectorsCmd() *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart connectors in a cluster",
		Long:  "",
		Run:   doRestartConnectors,
	}

	restartCmd.Flags().StringArrayVarP(&connectors, "connectors", "n", []string{}, "The connect names to restart (if not specified all connectors will be restarted)")
	_ = viper.BindPFlag("connectors", restartCmd.PersistentFlags().Lookup("connectors"))

	return restartCmd
}

func doRestartConnectors(cmd *cobra.Command, args []string) {
	clusterLogger := log.WithField("cluster", clusterURL)
	clusterLogger.Debug("executing manage connectors command")

	config := &manager.Config{
		ClusterURL: clusterURL,
		Logger:     clusterLogger,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("restart connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	err = mngr.Restart(connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatal("error restarting connectors")
	}

	clusterLogger.Info("connectors restarted successfully")
}
