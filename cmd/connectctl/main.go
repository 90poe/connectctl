package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/90poe/connectctl/pkg/ctl/manage"
	"github.com/90poe/connectctl/pkg/version"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "connectctl [command]",
		Short: "A kafka connect CLI",
		Long:  "",
		Run: func(c *cobra.Command, _ []string) {
			_ = c.Help()
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "c",
		"Config file (default is $HOME/.connectl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "warn", "Log level for the CLI")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.AddCommand(manage.Command())

	cobra.OnInitialize(initConfig)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).Errorf("error parsing log level: %s", logLevel)
	}
	log.SetLevel(level)

	log.Infof("connectctl, %s", version.ToString())

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
