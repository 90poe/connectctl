package connectors

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command creates the tehe management commands
func Command() *cobra.Command {
	connectorsCmd := &cobra.Command{
		Use:   "connectors",
		Short: "Commands related to Kafka Connect Connectors",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.WithError(err).Errorln("Error showing help")
			}
		},
	}

	// Add subcommands
	connectorsCmd.AddCommand(manageConnectorsCmd())
	connectorsCmd.AddCommand(restartConnectorsCmd())
	connectorsCmd.AddCommand(listConnectorsCmd())
	connectorsCmd.AddCommand(addConnectorCmd())
	connectorsCmd.AddCommand(removeConnectorCmd())
	connectorsCmd.AddCommand(pauseConnectorsCmd())
	connectorsCmd.AddCommand(resumeConnectorsCmd())

	return connectorsCmd
}
