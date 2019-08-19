package plugins

import (
	"encoding/json"
	"os"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type listPluginsCmdParams struct {
	ClusterURL string
	Output     string
}

func listPluginsCmd() *cobra.Command {
	params := &listPluginsCmdParams{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List connector plugins in a cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doListPlugins(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(listCmd, &params.ClusterURL)
	ctl.AddOutputFlags(listCmd, &params.Output)

	return listCmd
}

func doListPlugins(_ *cobra.Command, params *listPluginsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Debug("listing connector plugins")

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("list connector plugins configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	plugins, err := mngr.GetAllPlugins()
	if err != nil {
		clusterLogger.WithError(err).Fatal("error getting all connector plguns")
	}

	switch params.Output {
	case "json":
		printPluginsAsJSON(plugins, clusterLogger)
	case "table":
		printPluginsAsTable(plugins, clusterLogger)
	default:
		clusterLogger.Errorf("invalid output format specified: %s", params.Output)
	}
}

func printPluginsAsJSON(plugins []*connect.Plugin, logger *log.Entry) {
	logger.Debug("printing plugins as JSON")
	b, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		logger.WithError(err).Fatalf("error printing plugins as JSON")
	}

	os.Stdout.Write(b)
}

func printPluginsAsTable(plugins []*connect.Plugin, logger *log.Entry) {
	logger.Debug("printing plugins as table")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Class", "Type", "Version"})

	for _, plugin := range plugins {

		t.AppendRow(table.Row{
			plugin.Class,
			plugin.Type,
			plugin.Version,
		})
	}
	t.Render()
}
