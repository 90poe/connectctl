package connectors

import (
	"fmt"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doRestartConnectors(cmd, params)
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

func doRestartConnectors(_ *cobra.Command, params *restartConnectorsCmdParams) error {
	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}

	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(params.ClusterURL, connect.WithUserAgent(userAgent))
	if err != nil {
		return errors.Wrap(err, "error creating connect client")
	}

	mngr, err := manager.NewConnectorsManager(client, config)
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}

	if err = mngr.Restart(params.Connectors, params.RestartTasks, params.ForceRestartTasks, params.TaskIDs); err != nil {
		return errors.Wrap(err, "error restarting connectors")
	}

	fmt.Printf("connectors restarted successfully")
	return nil
}
