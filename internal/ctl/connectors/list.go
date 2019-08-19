package connectors

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
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
		Run: func(cmd *cobra.Command, _ []string) {
			doListConnectors(cmd, params)
		},
	}

	addCommonConnectorsFlags(listCmd, &params.ClusterURL)
	addOutputFlags(listCmd, &params.Output)

	return listCmd
}

func doListConnectors(_ *cobra.Command, params *listConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Debug("listing connectors")

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("list connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	connectors, err := mngr.GetAllConnectors()
	if err != nil {
		clusterLogger.WithError(err).Fatal("error getting all connectors")
	}

	switch params.Output {
	case "json":
		printConnectorsAsJSON(connectors, clusterLogger)
	case "table":
		printConnectorsAsTable(connectors, clusterLogger)
	default:
		clusterLogger.Errorf("invalid output format specified: %s", params.Output)
	}
}

func printConnectorsAsJSON(connectors []*manager.ConnectorWithState, logger *log.Entry) {
	logger.Debug("printing connectors as JSON")
	b, err := json.MarshalIndent(connectors, "", "  ")
	if err != nil {
		logger.WithError(err).Fatalf("error printing connectors as JSON")
	}

	os.Stdout.Write(b)
}

func printConnectorsAsTable(connectors []*manager.ConnectorWithState, logger *log.Entry) {
	logger.Debug("printing connectors as table")

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