package plugins

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command creates the plugins command (and subcommands)
func Command() *cobra.Command {
	pluginsCmd := &cobra.Command{
		Use:   "plugins",
		Short: "Commands related to Kafka Connect plugins",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.WithError(err).Errorln("Error showing help")
			}
		},
	}

	// Add subcommands
	pluginsCmd.AddCommand(listPluginsCmd())

	return pluginsCmd
}
