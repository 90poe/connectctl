package connectors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/healthcheck"
	"github.com/90poe/connectctl/internal/logging"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
	signals "github.com/90poe/connectctl/pkg/signal"
	"github.com/90poe/connectctl/pkg/sources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/matryer/try.v1"

	"github.com/kelseyhightower/envconfig"
)

type manageDefaults struct {
	ClusterURL                   string        `envconfig:"CLUSTER"`
	Files                        []string      `envconfig:"FILES"`
	Directory                    string        `envconfig:"DIRECTORY"`
	EnvVar                       string        `envconfig:"ENV_VAR"`
	InitialWaitPeriod            time.Duration `envconfig:"INITIAL_WAIT_PERIOD"`
	SyncPeriod                   time.Duration `envconfig:"SYNC_PERIOD"`
	SyncErrorRetryMax            int           `envconfig:"SYNC_ERROR_RETRY_MAX"`
	SyncErrorRetryPeriod         time.Duration `envconfig:"SYNC_ERROR_RETRY_PERIOD"`
	AllowPurge                   bool          `envconfig:"ALLOW_PURGE"`
	AutoRestart                  bool          `envconfig:"AUTO_RESTART"`
	RunOnce                      bool          `envconfig:"RUN_ONCE"`
	EnableHealthCheck            bool          `envconfig:"HEALTHCHECK_ENABLE"`
	HealthCheckAddress           string        `envconfig:"HEALTHCHECK_ADDRESS"`
	HTTPClientTimeout            time.Duration `envconfig:"HTTP_CLIENT_TIMEOUT"`
	GlobalConnectorRestartsMax   int           `envconfig:"GLOBAL_CONNECTOR_RESTARTS_MAX"`
	GlobalConnectorRestartPeriod time.Duration `envconfig:"GLOBAL_CONNECTOR_RESTART_PERIOD"`
	GlobalTaskRestartsMax        int           `envconfig:"GLOBAL_TASK_RESTARTS_MAX"`
	GlobalTaskRestartPeriod      time.Duration `envconfig:"GLOBAL_TASK_RESTART_PERIOD"`
	LogLevel                     string        `envconfig:"LOG_LEVEL"`
	LogFile                      string        `envconfig:"LOG_FILE"`
	LogFormat                    string        `envconfig:"LOG_FORMAT"`
}

func manageConnectorsCmd() *cobra.Command { // nolint: funlen
	params := &manageDefaults{
		InitialWaitPeriod:            3 * time.Minute,
		SyncPeriod:                   5 * time.Minute,
		SyncErrorRetryMax:            10,
		SyncErrorRetryPeriod:         1 * time.Minute,
		HealthCheckAddress:           ":9000",
		HTTPClientTimeout:            20 * time.Second,
		GlobalConnectorRestartsMax:   5,
		GlobalConnectorRestartPeriod: 10 * time.Second,
		GlobalTaskRestartsMax:        5,
		GlobalTaskRestartPeriod:      10 * time.Second,
		LogLevel:                     "INFO",
		LogFormat:                    "TEXT",
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
			err := envconfig.Process("CONNECTCTL", params)

			if err != nil {
				return errors.Wrap(err, "error processing environmental configuration")
			}

			return doManageConnectors(cmd, params)
		},
	}

	ctl.AddCommonConnectorsFlags(manageCmd, &params.ClusterURL)

	ctl.AddDefinitionFilesFlags(manageCmd, &params.Files, &params.Directory, &params.EnvVar)

	ctl.BindDurationVarP(manageCmd.Flags(), &params.SyncPeriod, params.SyncPeriod, "sync-period", "s", "how often to sync with the connect cluster")
	ctl.BindDurationVar(manageCmd.Flags(), &params.InitialWaitPeriod, params.InitialWaitPeriod, "wait-period", "time period to wait before starting the first sync")

	ctl.BindBoolVar(manageCmd.Flags(), &params.AllowPurge, false, "allow-purge", "if set connectctl will manage all connectors in a cluster. If connectors exist in the cluster that aren' t specified in --files then the connectors will be deleted")
	ctl.BindBoolVar(manageCmd.Flags(), &params.AutoRestart, false, "auto-restart", "if set connectors and tasks that are failed with automatically be restarted")
	ctl.BindBoolVar(manageCmd.Flags(), &params.RunOnce, false, "once", "if supplied sync will run once and command will exit")

	ctl.BindBoolVar(manageCmd.Flags(), &params.EnableHealthCheck, false, "healthcheck-enable", "if true a healthcheck via http will be enabled")
	ctl.BindStringVar(manageCmd.Flags(), &params.HealthCheckAddress, params.HealthCheckAddress, "healthcheck-address", "if enabled the healthchecks ('/live' and '/ready') will be available from this address")

	ctl.BindDurationVar(manageCmd.Flags(), &params.HTTPClientTimeout, params.HTTPClientTimeout, "http-client-timeout", "HTTP client timeout")

	ctl.BindIntVar(manageCmd.Flags(), &params.GlobalConnectorRestartsMax, params.GlobalConnectorRestartsMax, "global-connector-restarts-max", "maximum times a failed connector will be restarted")
	ctl.BindDurationVar(manageCmd.Flags(), &params.GlobalConnectorRestartPeriod, params.GlobalConnectorRestartPeriod, "global-connector-restart-period", "period of time between failed connector restart attemots")
	ctl.BindIntVar(manageCmd.Flags(), &params.GlobalTaskRestartsMax, params.GlobalTaskRestartsMax, "global-task-restarts-max", "maximum times a failed task will be restarted")
	ctl.BindDurationVar(manageCmd.Flags(), &params.GlobalTaskRestartPeriod, params.GlobalTaskRestartPeriod, "global-task-restart-period", "period of time between failed task restarts")

	ctl.BindIntVar(manageCmd.Flags(), &params.SyncErrorRetryMax, params.SyncErrorRetryMax, "sync-error-retry-max", "maximum times to ignore retryable errors whilst syncing")
	ctl.BindDurationVar(manageCmd.Flags(), &params.SyncErrorRetryPeriod, params.SyncErrorRetryPeriod, "sync-error-retry-period", "period of time between retryable errors whilst syncing")

	ctl.BindStringVarP(manageCmd.Flags(), &params.LogLevel, params.LogLevel, "loglevel", "l", "Log level for the CLI (Optional)")
	ctl.BindStringVar(manageCmd.Flags(), &params.LogFile, params.LogFile, "logfile", "A file to use for log output (Optional)")
	ctl.BindStringVar(manageCmd.Flags(), &params.LogFormat, params.LogFormat, "logformat", "Format for log output (Optional)")

	return manageCmd
}

func doManageConnectors(cmd *cobra.Command, params *manageDefaults) error {
	if err := logging.Configure(params.LogLevel, params.LogFile, params.LogFormat); err != nil {
		return errors.Wrap(err, "error configuring logging")
	}

	logger := log.WithFields(log.Fields{
		"cluster": params.ClusterURL,
		"version": version.Version,
	})
	logger.Debug("executing manage connectors command")

	if err := checkConfigSwitches(params.Files, params.Directory, params.EnvVar); err != nil {
		return errors.Wrap(err, "error with configuration")
	}

	config := &manager.Config{
		ClusterURL:                   params.ClusterURL,
		InitialWaitPeriod:            params.InitialWaitPeriod,
		SyncPeriod:                   params.SyncPeriod,
		AllowPurge:                   params.AllowPurge,
		AutoRestart:                  params.AutoRestart,
		Version:                      version.Version,
		GlobalConnectorRestartsMax:   params.GlobalConnectorRestartsMax,
		GlobalConnectorRestartPeriod: params.GlobalConnectorRestartPeriod,
		GlobalTaskRestartsMax:        params.GlobalTaskRestartsMax,
		GlobalTaskRestartPeriod:      params.GlobalTaskRestartPeriod,
	}

	logger.WithField("config", config).Trace("manage connectors configuration")

	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(params.ClusterURL, connect.WithUserAgent(userAgent), connect.WithHTTPClient(&http.Client{Timeout: params.HTTPClientTimeout}))
	if err != nil {
		return errors.Wrap(err, "error creating connect client")
	}

	mngr, err := manager.NewConnectorsManager(client, config, manager.WithLogger(logger))
	if err != nil {
		return errors.Wrap(err, "error creating connectors manager")
	}

	return syncOrManage(logger, params, cmd, mngr)
}

func syncOrManage(logger *log.Entry, params *manageDefaults, cmd *cobra.Command, mngr *manager.ConnectorManager) error {
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
				return true, errors.New("retry please")
			}
			lgr.Error("non-recoverable error when running manage")
			return false, ierr
		}
		lgr.Info("attempt finished")
		return false, nil
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

func checkConfigSwitches(files []string, directory, envar string) error {
	paramsSet := 0

	if len(files) != 0 {
		paramsSet++
	}
	if directory != "" {
		paramsSet++
	}
	if envar != "" {
		paramsSet++
	}

	if paramsSet == 1 {
		return nil
	}

	return errors.New("you must supply a list of files using --files or a directory that contains files using --directory or an environmental whose value is a JSON serialised connector or array of connectors")
}
