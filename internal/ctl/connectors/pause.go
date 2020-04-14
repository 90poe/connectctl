package connectors

import (
	"fmt"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

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

	if err = mngr.Pause(params.Connectors); err != nil {
		return errors.Wrap(err, "error pausing connectors")
	}

	fmt.Println("connectors paused successfully")
	return nil
}
