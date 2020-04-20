package connectors

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type connectorsStatusCmdParams struct {
	ClusterURL string
	Connectors []string
	Output     string
	Quiet      bool
}

func connectorsStatusCmd() *cobra.Command {
	params := &connectorsStatusCmdParams{}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get status for connectors in a cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doConnectorsStatus(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(statusCmd, &params.ClusterURL)
	ctl.AddConnectorNamesFlags(statusCmd, &params.Connectors)
	ctl.AddOutputFlags(statusCmd, &params.Output)
	ctl.AddQuietFlag(statusCmd, &params.Quiet)

	return statusCmd
}

func doConnectorsStatus(_ *cobra.Command, params *connectorsStatusCmdParams) error {
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

	statusList, err := mngr.Status(params.Connectors)
	if err != nil {
		return errors.Wrap(err, "error getting connectors status")
	}

	if !params.Quiet {
		switch params.Output {
		case "json":
			if err = printAsJSON(statusList); err != nil {
				return errors.Wrap(err, "error printing connectors status as JSON")
			}

		case "table":
			printAsTable(statusList)

		default:
			return fmt.Errorf("invalid output format specified: %s", params.Output)
		}
	}

	failingConnectors, failingTasks := countFailing(statusList)

	if failingConnectors != 0 || failingTasks != 0 {
		return fmt.Errorf("%d connectors are failng, %d tasks are failing", failingConnectors, failingTasks)
	}

	return nil
}

func countFailing(statusList []*connect.ConnectorStatus) (int, int) {
	connectorCount := 0
	taskCount := 0

	for _, status := range statusList {
		if status.Connector.State == "FAILED" {
			connectorCount++
		}

		taskCount += countFailingTasks(&status.Tasks)
	}

	return connectorCount, taskCount
}

func countFailingTasks(tasks *[]connect.TaskState) int {
	count := 0

	for _, task := range *tasks {
		if task.State == "FAILED" {
			count++
		}
	}

	return count
}

func printAsJSON(statusList []*connect.ConnectorStatus) error {
	b, err := json.MarshalIndent(statusList, "", "  ")
	if err != nil {
		return err
	}

	os.Stdout.Write(b)
	return nil
}

func printAsTable(statusList []*connect.ConnectorStatus) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "State", "WorkerId", "Tasks"})

	for _, status := range statusList {
		tasks := ""
		for _, task := range status.Tasks {
			tasks += fmt.Sprintf("%d(%s): %s\n", task.ID, task.WorkerID, task.State)
		}

		t.AppendRow(table.Row{
			status.Name,
			status.Connector.State,
			status.Connector.WorkerID,
			tasks,
		})
	}

	t.Render()
}
