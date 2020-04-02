package tasks

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type RestartOptions struct {
	*GenericOptions

	TaskID int
}

func NewRestartCommand(options *GenericOptions) *cobra.Command {
	var opts = RestartOptions{
		GenericOptions: options,
	}

	var cmd = cobra.Command{
		Use:   "restart",
		Short: "Restart an individual task for the specified connector.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return opts.Run()
		},
	}

	cmd.PersistentFlags().IntVarP(&opts.TaskID, "id", "", -1, "The ID of the task to restart")
	_ = cmd.MarkPersistentFlagRequired("id")

	return &cmd
}

func (o *RestartOptions) Run() error {
	client, err := o.CreateClient(o.ClusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create http client")
	}
	_, err = client.RestartConnectorTask(o.ConnectorName, o.TaskID)
	if err != nil {
		return errors.Wrapf(err, "failed to restart task '%d' for connector '%s'", o.TaskID, o.ConnectorName)
	}
	return nil
}
