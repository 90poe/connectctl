package connectors

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/pkg/errors"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type listConnectorsCmdParams struct {
	ClusterURL string
	Output     string
}

func listConnectorsCmd() *cobra.Command {
	params := &listConnectorsCmdParams{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List connectors in a cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doListConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(listCmd, &params.ClusterURL)
	ctl.AddOutputFlags(listCmd, &params.Output)

	return listCmd
}

func doListConnectors(_ *cobra.Command, params *listConnectorsCmdParams) error {

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

	connectors, err := mngr.GetAllConnectors()
	if err != nil {
		return errors.Wrap(err, "error getting all connectors")
	}

	switch params.Output {
	case "json":
		err := printConnectorsAsJSON(connectors)
		if err != nil {
			return errors.Wrap(err, "error printing connectors as JSON")
		}
	case "table":
		printConnectorsAsTable(connectors)
	default:
		return fmt.Errorf("invalid output format specified: %s", params.Output)
	}
	return nil
}

func printConnectorsAsJSON(connectors []*manager.ConnectorWithState) error {
	b, err := json.MarshalIndent(connectors, "", "  ")
	if err != nil {
		return err
	}

	os.Stdout.Write(b)
	return nil
}

func printConnectorsAsTable(connectors []*manager.ConnectorWithState) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "State", "WorkerId", "Tasks", "Config"})

	for _, connector := range connectors {
		config := ""
		for key, val := range connector.Config {
			config += fmt.Sprintf("%s=%s\n", key, val)
		}

		tasks := ""
		for _, task := range connector.Tasks {
			tasks += fmt.Sprintf("%d(%s): %s\n", task.ID, task.WorkerID, task.State)
		}

		t.AppendRow(table.Row{
			connector.Name,
			connector.ConnectorState.State,
			connector.ConnectorState.WorkerID,
			tasks,
			config,
		})
	}
	t.Render()
}
