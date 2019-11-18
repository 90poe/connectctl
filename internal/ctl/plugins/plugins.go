package plugins

import (
	"github.com/spf13/cobra"
)

// Command creates the plugins command (and subcommands)
func Command() *cobra.Command {
	pluginsCmd := &cobra.Command{
		Use:   "plugins",
		Short: "Commands related to Kafka Connect plugins",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	// Add subcommands
	pluginsCmd.AddCommand(listPluginsCmd())

	return pluginsCmd
}
