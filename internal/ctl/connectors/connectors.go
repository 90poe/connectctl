package connectors

import (
	"github.com/spf13/cobra"
)

// Command creates the the management commands
func Command() *cobra.Command {
	connectorsCmd := &cobra.Command{
		Use:   "connectors",
		Short: "Commands related to Kafka Connect Connectors",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
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
	connectorsCmd.AddCommand(connectorsStatusCmd())

	return connectorsCmd
}
