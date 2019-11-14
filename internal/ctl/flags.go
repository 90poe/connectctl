package ctl

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddCommonConnectorsFlags(cmd *cobra.Command, clusterURL *string) {
	cmd.PersistentFlags().StringVarP(clusterURL, "cluster", "c", "", "the URL of the connect cluster to manage (required)")
	_ = cmd.MarkPersistentFlagRequired("cluster")
	_ = viper.BindPFlag("cluster", cmd.PersistentFlags().Lookup("cluster"))
}

func AddOutputFlags(cmd *cobra.Command, output *string) {
	cmd.PersistentFlags().StringVarP(output, "output", "o", "json", "specified the output format (valid options: json, table)")
	_ = viper.BindPFlag("output", cmd.PersistentFlags().Lookup("output"))
}

func AddDefinitionFilesFlags(cmd *cobra.Command, files *[]string, directory *string, env *string) {
	cmd.Flags().StringArrayVarP(files, "files", "f", []string{}, "The connector definitions files (Required if --directory or --env-var not specified)")
	_ = viper.BindPFlag("files", cmd.PersistentFlags().Lookup("files"))

	cmd.Flags().StringVarP(directory, "directory", "d", "", "The directory containing the connector definitions files (Required if --file or --env-vars not specified)")
	_ = viper.BindPFlag("directory", cmd.PersistentFlags().Lookup("directory"))

	cmd.Flags().StringVarP(env, "env-var", "e", "", "An environmental variable whose value is a singular or array of connectors serialised as JSON (Required if --files or --directory not specified)")
	_ = viper.BindPFlag("env-var", cmd.PersistentFlags().Lookup("env-var"))

}

func AddConnectorNamesFlags(cmd *cobra.Command, names *[]string) {
	cmd.Flags().StringArrayVarP(names, "connectors", "n", []string{}, "The connect names to restart (if not specified all connectors will be restarted)")
	_ = viper.BindPFlag("connectors", cmd.PersistentFlags().Lookup("connectors"))
}
