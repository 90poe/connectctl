package tasks

import (
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strings"
)

type ClientFn func(string) (Client, error)

type Client interface {
	GetConnectorTasks(name string) ([]connect.Task, *http.Response, error)
	GetConnectorTaskStatus(name string, taskID int) (*connect.TaskState, *http.Response, error)
	RestartConnectorTask(name string, taskID int) (*http.Response, error)
}

type GenericOptions struct {
	// Function for creating API client
	CreateClient ClientFn

	ClusterURL    string
	ConnectorName string
}

func Command(opts *GenericOptions) *cobra.Command {
	var cmd = cobra.Command{
		Use:   "tasks",
		Short: "Commands related to kafka connector tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			// return 'help' by default
			return cmd.Help()
		},
	}

	// ClusterURL will be used in sub-commands
	cmd.PersistentFlags().StringVarP(&opts.ClusterURL, "cluster", "c", "", "the URL of the connect cluster to manage (required)")
	_ = cmd.MarkPersistentFlagRequired("cluster")

	// ConnectorName will be used in sub-commands
	cmd.PersistentFlags().StringVar(&opts.ConnectorName, "connector", "", "Connector name to get tasks for")
	_ = cmd.MarkPersistentFlagRequired("connector")

	// connectctl tasks list --cluster=... --connector=...
	cmd.AddCommand(NewListCommand(opts))

	// connectctl task get --cluster=... --connector=... --id=...
	cmd.AddCommand(NewGetCommand(opts))

	// connectctl task restart --cluster=... --connector=... --id=...
	cmd.AddCommand(NewRestartCommand(opts))

	// connectctl task status --cluster=... --connector=... --id=...
	cmd.AddCommand(NewStatusCommand(opts))

	return &cmd
}

func (o *GenericOptions) Validate() error {
	_, err := url.ParseRequestURI(o.ClusterURL)
	if err != nil {
		return errors.Wrap(err, "--cluster is not a valid URI")
	}
	if len(strings.TrimSpace(o.ConnectorName)) == 0 {
		return errors.New("--connector name is empty")
	}
	return nil
}
