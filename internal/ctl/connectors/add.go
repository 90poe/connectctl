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

type addConnectorsCmdParams struct {
	ClusterURL string
	Files      []string
	Directory  string
	EnvVar     string
}

func addConnectorCmd() *cobra.Command {
	params := &addConnectorsCmdParams{}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add connectors to a connect cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doAddConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(addCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(addCmd, &params.Files, &params.Directory, &params.EnvVar)

	return addCmd
}

func doAddConnectors(cmd *cobra.Command, params *addConnectorsCmdParams) error {

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}

	source, err := findSource(params.Files, params.Directory, params.EnvVar, cmd)

	if err != nil {
		return err
	}

	connectors, err := source()
	if err != nil {
		return errors.Wrap(err, "error reading connector configuration from files")
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
	if err = mngr.Add(connectors); err != nil {
		return errors.Wrap(err, "error creating connectors")
	}
	fmt.Println("added connectors")
	return nil
}
