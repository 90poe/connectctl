package restart

import (
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clusterURL string
	connectors []string
)

// Command creates the tehe management commands
func Command() *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart connectors in a cluster",
		Long:  "",
		Run:   doRestartConnectors,
	}

	// Add persistent flags that apply to all sub commands
	restartCmd.PersistentFlags().StringVarP(&clusterURL, "cluster", "c", "", "The URL of the connect cluster to manage (required)")
	_ = restartCmd.MarkPersistentFlagRequired("cluster")
	_ = viper.BindPFlag("cluster", restartCmd.PersistentFlags().Lookup("cluster"))

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
	}
	clusterLogger.WithField("config", config).Trace("restart connectors confirguration")

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
