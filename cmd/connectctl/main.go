package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/90poe/connectctl/internal/ctl/manage"
	"github.com/90poe/connectctl/internal/ctl/restart"
	"github.com/90poe/connectctl/internal/logging"
	"github.com/90poe/connectctl/internal/version"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
	logFile  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "connectctl [command]",
		Short: "A kafka connect CLI",
		Long:  "",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			err := logging.Configure(logLevel, logFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			log.Infof("connectctl, %s", version.ToString())
		},
		Run: func(c *cobra.Command, _ []string) {
			_ = c.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "Config file (default is $HOME/.connectl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "", "Log level for the CLI (Optional)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "logfile", "", "", "A file to use for log output (Optional)")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	_ = viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.SetDefault("loglevel", "INFO")

	rootCmd.AddCommand(manage.Command())
	rootCmd.AddCommand(restart.Command())

	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
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
