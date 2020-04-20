package plugins

import (
	"fmt"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type validatePluginsCmdParams struct {
	ClusterURL string
	Input      string
}

func validatePluginsCmd() *cobra.Command {
	params := &validatePluginsCmdParams{}

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates plugin config",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doValidatePlugins(cmd, params)
		},
	}

	ctl.AddClusterFlag(validateCmd, true, &params.ClusterURL)
	ctl.AddInputFlag(validateCmd, true, &params.Input)

	return validateCmd
}

func doValidatePlugins(_ *cobra.Command, params *validatePluginsCmdParams) error {
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

	//TODO remove
	if mngr != nil {
		return nil
	}

	return nil
}
