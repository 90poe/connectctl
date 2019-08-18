package connectors

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clusterURL string
)

// Command creates the tehe management commands
func Command() *cobra.Command {
	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Actively manage a connect cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.WithError(err).Errorln("Error showing help")
			}
		},
	}

	// Add persistent flags that apply to all sub commands
	manageCmd.PersistentFlags().StringVarP(&clusterURL, "cluster", "c", "", "The URL of the connect cluster to manage (required)")
	_ = manageCmd.MarkPersistentFlagRequired("cluster")
	_ = viper.BindPFlag("cluster", manageCmd.PersistentFlags().Lookup("cluster"))

	// Add subcommands
	manageCmd.AddCommand(manageConnectorsCmd())
	manageCmd.AddCommand(restartConnectorsCmd())

	return manageCmd
}
