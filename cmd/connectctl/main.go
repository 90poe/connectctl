//nolint:gochecknoglobals
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/90poe/connectctl/internal/ctl/tasks"
	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/90poe/connectctl/internal/ctl/connectors"
	"github.com/90poe/connectctl/internal/ctl/plugins"
	"github.com/90poe/connectctl/internal/ctl/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "connectctl [command]",
		Short: "A kafka connect CLI",
		Long:  "",
		RunE: func(c *cobra.Command, _ []string) error {
			return c.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "Config file (default is $HOME/.connectctl.yaml)")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(connectors.Command())
	rootCmd.AddCommand(plugins.Command())
	rootCmd.AddCommand(version.Command())

	rootCmd.AddCommand(tasks.Command(&tasks.GenericOptions{
		CreateClient: func(clusterURL string) (client tasks.Client, err error) {
			return connect.NewClient(clusterURL)
		},
	}))

	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigFile(".connectctl.yaml")
	}

	replacer := strings.NewReplacer(".", "-")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
