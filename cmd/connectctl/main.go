//nolint:gochecknoglobals
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/90poe/connectctl/internal/ctl/connectors"
	"github.com/90poe/connectctl/internal/ctl/plugins"
	"github.com/90poe/connectctl/internal/ctl/version"
	"github.com/90poe/connectctl/internal/logging"

	"github.com/pkg/errors"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	logLevel  string
	logFile   string
	logFormat string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "connectctl [command]",
		Short: "A kafka connect CLI",
		Long:  "",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			err := logging.Configure(logLevel, logFile, logFormat)
			if err != nil {
				return errors.Wrap(err, "error configuring logging")
			}
			log.Info("connectctl, a Kafka Connect CLI")
			return nil
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return c.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "Config file (default is $HOME/.connectctl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "INFO", "Log level for the CLI (Optional)")
	rootCmd.PersistentFlags().StringVarP(&logFile, "logfile", "", "", "A file to use for log output (Optional)")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "", "TEXT", "Format for log output (Optional)")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	_ = viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	_ = viper.BindPFlag("logoutput", rootCmd.PersistentFlags().Lookup("logformat"))

	viper.SetDefault("loglevel", "INFO")

	rootCmd.AddCommand(connectors.Command())
	rootCmd.AddCommand(plugins.Command())
	rootCmd.AddCommand(version.Command())

	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v", err)
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
