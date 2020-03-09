package connectors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/healthcheck"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	signals "github.com/90poe/connectctl/pkg/signal"
	"github.com/90poe/connectctl/pkg/sources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/matryer/try.v1"
)

type manageConnectorsCmdParams struct {
	ClusterURL                   string
	Files                        []string
	Directory                    string
	EnvVar                       string
	SyncPeriod                   time.Duration
	SyncErrorRetryMax            int
	SyncErrorRetryPeriod         time.Duration
	AllowPurge                   bool
	AutoRestart                  bool
	RunOnce                      bool
	EnableHealthCheck            bool
	HealthCheckAddress           string
	HTTPClientTimeout            time.Duration
	GlobalMaxConnectorRestarts   int
	GlobalConnectorRestartPeriod time.Duration
	GlobalMaxTaskRestarts        int
	GlobalTaskRestartPeriod      time.Duration
}

func manageConnectorsCmd() *cobra.Command { // nolint: funlen
	params := &manageConnectorsCmdParams{
		SyncPeriod:                   5 * time.Minute,
		SyncErrorRetryMax:            10,
		SyncErrorRetryPeriod:         1 * time.Minute,
		HealthCheckAddress:           ":9000",
		HTTPClientTimeout:            20 * time.Second,
		GlobalMaxConnectorRestarts:   5,
		GlobalConnectorRestartPeriod: 10 * time.Second,
		GlobalMaxTaskRestarts:        5,
		GlobalTaskRestartPeriod:      10 * time.Second,
	}

	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Actively manage connectors in a Kafka Connect cluster",
		Long: `This command will add/delete/update connectors in a destination
Kafa Connect cluster based on a list of desired connectors which are specified
as a list of files or all files in a directory. The command runs continuously and
will sync desired state with actual state based on the --sync-period flag. But
if you specify --once then it will sync once and then exit.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doManageConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(manageCmd, &params.ClusterURL)
	ctl.AddDefinitionFilesFlags(manageCmd, &params.Files, &params.Directory, &params.EnvVar)

	manageCmd.Flags().DurationVarP(&params.SyncPeriod, "sync-period", "s", params.SyncPeriod, "how often to sync with the connect cluster")
	_ = viper.BindPFlag("sync-period", manageCmd.PersistentFlags().Lookup("sync-period"))

	manageCmd.Flags().BoolVarP(&params.AllowPurge, "allow-purge", "", false, "if set connectctl will manage all connectors in a cluster. If connectors exist in the cluster that aren't specified in --files then the connectors will be deleted")
	_ = viper.BindPFlag("allow-purge", manageCmd.PersistentFlags().Lookup("allow-purge"))

	manageCmd.Flags().BoolVar(&params.AutoRestart, "auto-restart", false, "if set connectors and tasks that are failed with automatically be restarted")
	_ = viper.BindPFlag("auto-restart", manageCmd.PersistentFlags().Lookup("auto-restart"))

	manageCmd.Flags().BoolVar(&params.RunOnce, "once", false, "if supplied sync will run once and command will exit")
	_ = viper.BindPFlag("once", manageCmd.PersistentFlags().Lookup("once"))

	manageCmd.Flags().BoolVar(&params.EnableHealthCheck, "healthcheck-enable", false, "if supplied a healthcheck via http will be enabled")
	_ = viper.BindPFlag("healthcheck-enable", manageCmd.PersistentFlags().Lookup("healthcheck-enable"))

	manageCmd.Flags().StringVar(&params.HealthCheckAddress, "healthcheck-address", params.HealthCheckAddress, "if enabled the healthchecks ('/live' and '/ready') will be available from this address")
	_ = viper.BindPFlag("healthcheck-address", manageCmd.PersistentFlags().Lookup("healthcheck-address"))

	manageCmd.Flags().DurationVar(&params.HTTPClientTimeout, "http-client-timeout", params.HTTPClientTimeout, "HTTP client timeout")
	_ = viper.BindPFlag("http-client-timeout", manageCmd.PersistentFlags().Lookup("http-client-timeout"))

	manageCmd.Flags().IntVar(&params.GlobalMaxConnectorRestarts, "global-connector-max-restarts", params.GlobalMaxConnectorRestarts, "maximum times a failed connector will be restarted")
	_ = viper.BindPFlag("global-connector-max-restarts", manageCmd.PersistentFlags().Lookup("global-connector-max-restarts"))

	manageCmd.Flags().DurationVar(&params.GlobalConnectorRestartPeriod, "global-connector-restart-period", params.GlobalConnectorRestartPeriod, "period of time between failed connector restarts")
	_ = viper.BindPFlag("global-connector-restart-period", manageCmd.PersistentFlags().Lookup("global-connector-restart-period"))

	manageCmd.Flags().IntVar(&params.GlobalMaxTaskRestarts, "global-task-max-restarts", params.GlobalMaxTaskRestarts, "maximum times a failed task will be restarted")
	_ = viper.BindPFlag("global-task-max-restarts", manageCmd.PersistentFlags().Lookup("global-task-max-restarts"))

	manageCmd.Flags().DurationVar(&params.GlobalTaskRestartPeriod, "global-task-restart-period", params.GlobalTaskRestartPeriod, "period of time between failed task restarts")
	_ = viper.BindPFlag("global-task-restart-period", manageCmd.PersistentFlags().Lookup("global-task-restart-period"))

	return manageCmd
}

func doManageConnectors(cmd *cobra.Command, params *manageConnectorsCmdParams) error {
	logger := log.WithField("cluster", params.ClusterURL)
	logger.Debug("executing manage connectors command")

	if err := checkConfig(params); err != nil {
		return errors.Wrap(err, "Error with configuration")
	}

	config := &manager.Config{
		ClusterURL:                   params.ClusterURL,
		SyncPeriod:                   params.SyncPeriod,
		AllowPurge:                   params.AllowPurge,
		AutoRestart:                  params.AutoRestart,
		Version:                      version.Version,
		GlobalMaxConnectorRestarts:   params.GlobalMaxConnectorRestarts,
		GlobalConnectorRestartPeriod: params.GlobalConnectorRestartPeriod,
		GlobalMaxTaskRestarts:        params.GlobalMaxTaskRestarts,
		GlobalTaskRestartPeriod:      params.GlobalTaskRestartPeriod,
	}

	logger.WithField("config", config).Trace("manage connectors configuration")

	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(params.ClusterURL, connect.WithUserAgent(userAgent), connect.WithHTTPClient(&http.Client{Timeout: params.HTTPClientTimeout}))
	if err != nil {
		return errors.Wrap(err, "error creating connect client")
	}

	mngr, err := manager.NewConnectorsManager(client, config)
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}

	return syncOrManage(logger, params, cmd, mngr)
}

func syncOrManage(logger *log.Entry, params *manageConnectorsCmdParams, cmd *cobra.Command, mngr *manager.ConnectorManager) error {
	if params.EnableHealthCheck {
		healthCheckHandler := healthcheck.New(mngr)

		go func() {
			if err := healthCheckHandler.Start(params.HealthCheckAddress); err != nil {
				logger.WithError(err).Fatalln("error starting healthcheck")
			}
		}()
		// nolint
		defer healthCheckHandler.Shutdown(context.Background())
	}

	source, err := findSource(params.Files, params.Directory, params.EnvVar, cmd)

	if err != nil {
		return err
	}

	stopCh := signals.SetupSignalHandler()

	try.MaxRetries = params.SyncErrorRetryMax

	return try.Do(func(attempt int) (bool, error) {
		lgr := logger.WithField("attempt", attempt)

		var ierr error
		if params.RunOnce {
			lgr.Info("running once")
			ierr = mngr.Sync(source)
		} else {
			lgr.Info("managing")
			ierr = mngr.Manage(source, stopCh)
		}

		if ierr != nil {
			lgr = logger.WithError(ierr)
			rootCause := errors.Cause(ierr)
			if connect.IsRetryable(rootCause) {
				lgr.WithField("attempt", attempt).Error("recoverable error when running manage")
				time.Sleep(params.SyncErrorRetryPeriod)
			} else {
				lgr.Error("non-recoverable error when running manage")
				return false, ierr
			}
		}
		lgr.Info("attempt finished")
		return true, errors.New("retry please")
	})
}
func findSource(files []string, directory, envar string, cmd *cobra.Command) (manager.ConnectorSource, error) {
	switch {
	case len(files) > 0:
		if len(files) == 1 && files[0] == "-" {
			return sources.StdIn(cmd.InOrStdin()), nil
		}
		return sources.Files(files), nil

	case directory != "":
		return sources.Directory(directory), nil
	case envar != "":
		return sources.EnvVarValue(envar), nil
	}
	return nil, errors.New("error finding connector definitions from parameters")
}

func checkConfig(params *manageConnectorsCmdParams) error {
	paramsSet := 0

	if len(params.Files) != 0 {
		paramsSet++
	}
	if params.Directory != "" {
		paramsSet++
	}
	if params.EnvVar != "" {
		paramsSet++
	}

	if paramsSet == 1 {
		return nil
	}

	return errors.New("you must supply a list of files using --files or a directory that contains files using --directory or an environmental whose value is a JSON serialised connector or array of connectors")
}
