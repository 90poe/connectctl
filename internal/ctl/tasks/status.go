package tasks

import (
	"fmt"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

type StatusOptions struct {
	*GenericOptions

	TaskID int
	Output string
}

func NewStatusCommand(options *GenericOptions) *cobra.Command {
	var opts = StatusOptions{
		GenericOptions: options,
	}

	var cmd = cobra.Command{
		Use:   "status",
		Short: "Displays a status by individual task currently running for the connector.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return opts.Run(cmd.OutOrStdout())
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "json", "The output format. Possible values are 'json' and 'table'")

	cmd.PersistentFlags().IntVarP(&opts.TaskID, "id", "", -1, "The ID of the task to get status for")
	_ = cmd.MarkPersistentFlagRequired("id")

	return &cmd
}

func (o *StatusOptions) Run(out io.Writer) error {
	client, err := o.CreateClient(o.ClusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create http client")
	}
	taskState, _, err := client.GetConnectorTaskStatus(o.ConnectorName, o.TaskID)
	if err != nil {
		return errors.Wrap(err, "failed to get task status")
	}
	if taskState == nil {
		return errors.New("task state response is nil")
	}
	return o.writeOutput(taskState, out)
}

func (o *StatusOptions) writeOutput(taskState *connect.TaskState, out io.Writer) error {
	var outputFn TaskStateOutputFn
	var outputType = OutputType(o.Output)
	switch outputType {
	case OutputTypeJSON:
		outputFn = OutputTaskStateAsJSON
	case OutputTypeTable:
		outputFn = OutputTaskStateAsTable
	default:
		return fmt.Errorf("output type '%s' is not supported", o.Output)
	}
	output, err := outputFn(taskState)
	if err != nil {
		return fmt.Errorf("failed to form output for '%s' type", outputType)
	}
	if _, err := out.Write(output); err != nil {
		return errors.Wrap(err, "failed to write output")
	}
	return nil
}
