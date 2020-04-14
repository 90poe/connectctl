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

type resumeConnectorsCmdParams struct {
	ClusterURL string
	Connectors []string
}

func resumeConnectorsCmd() *cobra.Command {
	params := &resumeConnectorsCmdParams{}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume connectors in a cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doResumeConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(resumeCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(resumeCmd, &params.Connectors)

	return resumeCmd
}

func doResumeConnectors(_ *cobra.Command, params *resumeConnectorsCmdParams) error {
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

	if err = mngr.Resume(params.Connectors); err != nil {
		return errors.Wrap(err, "error resuming connectors")
	}

	fmt.Println("connectors resumed successfully")
	return nil
}
