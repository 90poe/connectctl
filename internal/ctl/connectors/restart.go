package connectors

import (
	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type restartConnectorsCmdParams struct {
	ClusterURL        string
	Connectors        []string
	ForceRestartTasks bool
	RestartTasks      bool
	TaskIDs           []int
}

// Command creates the the management commands
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

	restartCmd.Flags().BoolVar(&params.RestartTasks, "restart-tasks", true, "Whether to restart the connector tasks")
	_ = viper.BindPFlag("restart-tasks", restartCmd.PersistentFlags().Lookup("restart-tasks"))

	restartCmd.Flags().BoolVar(&params.ForceRestartTasks, "force-restart-tasks", false, "Whether to force restart the connector tasks")
	_ = viper.BindPFlag("force-restart-tasks", restartCmd.PersistentFlags().Lookup("force-restart-tasks"))

	restartCmd.Flags().IntSliceVarP(&params.TaskIDs, "tasks", "t", []int{}, "The task ids to restart (if no ids are specified, all connectors will be restarted)")
	_ = viper.BindPFlag("tasks", restartCmd.PersistentFlags().Lookup("tasks"))

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

	err = mngr.Restart(params.Connectors, params.RestartTasks, params.ForceRestartTasks, params.TaskIDs)
	if err != nil {
		clusterLogger.WithError(err).Fatal("error restarting connectors")
	}

	clusterLogger.Info("connectors restarted successfully")
}
