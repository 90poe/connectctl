package ctl

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func AddClusterFlag(cmd *cobra.Command, required bool, clusterURL *string) {
	description := "the URL of the connect cluster to manage"

	if required {
		description = requiredDescription(&description)
	}

	BindStringVarP(cmd.Flags(), clusterURL, "", "cluster", "c", description)
}

func AddCommonConnectorsFlags(cmd *cobra.Command, clusterURL *string) {
	AddClusterFlag(cmd, true, clusterURL)
}

func AddOutputFlags(cmd *cobra.Command, output *string) {
	BindStringVarP(cmd.Flags(), output, "json", "output", "o", "specify the output format (valid options: json, table)")
}

func AddQuietFlag(cmd *cobra.Command, quiet *bool) {
	BindBoolVarP(cmd.Flags(), quiet, false, "quiet", "q", "disable output logging")
}

func AddDefinitionFilesFlags(cmd *cobra.Command, files *[]string, directory *string, env *string) {
	BindStringArrayVarP(cmd.Flags(), files, []string{}, "files", "f", "the connector definitions files (Required if --directory or --env-var not specified)")
	BindStringVarP(cmd.Flags(), directory, "", "directory", "d", "the directory containing the connector definitions files (Required if --file or --env-vars not specified)")
	BindStringVarP(cmd.Flags(), env, "", "env-var", "e", "an environmental variable whose value is a singular or array of connectors serialised as JSON (Required if --files or --directory not specified)")
}

func AddConnectorNamesFlags(cmd *cobra.Command, names *[]string) {
	BindStringArrayVarP(cmd.Flags(), names, []string{}, "connectors", "n", "The connect names to restart (if not specified all connectors will be restarted)")
}

func BindDurationVarP(f *pflag.FlagSet, p *time.Duration, value time.Duration, long, short, description string) {
	f.DurationVarP(p, long, short, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindDurationVar(f *pflag.FlagSet, p *time.Duration, value time.Duration, long, description string) {
	f.DurationVar(p, long, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindBoolVar(f *pflag.FlagSet, p *bool, value bool, long, description string) {
	f.BoolVar(p, long, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindBoolVarP(f *pflag.FlagSet, p *bool, value bool, long, short, description string) {
	f.BoolVarP(p, long, short, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindStringVarP(f *pflag.FlagSet, p *string, value, long, short, description string) {
	f.StringVarP(p, long, short, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindStringVar(f *pflag.FlagSet, p *string, value, long, description string) {
	f.StringVar(p, long, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindStringArrayVarP(f *pflag.FlagSet, p *[]string, value []string, long, short, description string) {
	f.StringArrayVarP(p, long, short, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func BindIntVar(f *pflag.FlagSet, p *int, value int, long, description string) {
	f.IntVar(p, long, value, description)
	_ = viper.BindPFlag(long, f.Lookup(long))
	viper.SetDefault(long, value)
}

func requiredDescription(desc *string) string {
	return fmt.Sprintf("%s (required)", *desc)
}
