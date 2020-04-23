package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type validatePluginsCmdParams struct {
	ClusterURL string
	Input      string
	Output     string
	Quiet      bool
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
	ctl.AddOutputFlags(validateCmd, &params.Output)
	ctl.AddQuietFlag(validateCmd, &params.Quiet)

	return validateCmd
}

func doValidatePlugins(_ *cobra.Command, params *validatePluginsCmdParams) error {
	var inputConfig connect.ConnectorConfig
	if err := json.Unmarshal([]byte(params.Input), &inputConfig); err != nil {
		return errors.Wrap(err, "error parsing input connector config")
	}

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

	validation, err := mngr.ValidatePlugins(inputConfig)
	if err != nil {
		return err
	}

	if !params.Quiet {
		switch params.Output {
		case "json":
			if err = ctl.PrintAsJSON(validation); err != nil {
				return errors.Wrap(err, "error printing validation results as JSON")
			}

		case "table":
			printAsTable(validation)

		default:
			return fmt.Errorf("invalid output format specified: %s", params.Output)
		}
	}

	if validation.ErrorCount > 0 {
		return fmt.Errorf("detected %d errors in the configuation", validation.ErrorCount)
	}

	return nil
}

func printAsTable(validation *connect.ConfigValidation) {
	ctl.PrintAsTable(func(t table.Writer) {
		t.Style().Options.SeparateRows = true
		t.AppendHeader(table.Row{"Name", "Spec", "Value", "Errors"})

		for _, info := range validation.Configs {
			spec := fmt.Sprintf(
				"default: %s\nrequired: %v",
				ctl.StrPtrToStr(info.Definition.DefaultValue),
				info.Definition.Required,
			)

			errors := strings.Join(info.Value.Errors, "\n")

			t.AppendRow(table.Row{
				info.Definition.Name,
				spec,
				ctl.StrPtrToStr(info.Value.Value),
				errors,
			})
		}
	})
}
