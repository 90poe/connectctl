package tasks

import (
	"fmt"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"strconv"
)

type GetOptions struct {
	*GenericOptions

	TaskID int
	Output string
}

func NewGetCommand(options *GenericOptions) *cobra.Command {
	var opts = GetOptions{
		GenericOptions: options,
	}

	var cmd = cobra.Command{
		Use:   "get",
		Short: "Displays a single task currently running for the connector.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return opts.Run(cmd.OutOrStdout())
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "json", "The output format. Possible values are 'json' and 'table'")

	cmd.PersistentFlags().IntVarP(&opts.TaskID, "id", "", -1, "The ID of the task to get")
	_ = cmd.MarkPersistentFlagRequired("id")

	return &cmd
}

func (o *GetOptions) Run(out io.Writer) error {
	client, err := o.CreateClient(o.ClusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create http client")
	}
	tasks, _, err := client.GetConnectorTasks(o.ConnectorName)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve tasks")
	}
	task, ok := findTaskByID(tasks, o.TaskID)
	if !ok {
		return errors.New("no task found by id=" + strconv.Itoa(o.TaskID))
	}
	return o.writeOutput(task, out)
}

func (o *GetOptions) writeOutput(task connect.Task, out io.Writer) error {
	var outputFn TaskListOutputFn
	var outputType = OutputType(o.Output)
	switch outputType {
	case OutputTypeJSON:
		outputFn = OutputTaskListAsJSON
	case OutputTypeTable:
		outputFn = OutputTaskListAsTable
	default:
		return fmt.Errorf("output type '%s' is not supported", o.Output)
	}
	output, err := outputFn([]connect.Task{task})
	if err != nil {
		return fmt.Errorf("failed to form output for '%s' type", outputType)
	}
	if _, err := out.Write(output); err != nil {
		return errors.Wrap(err, "failed to write output")
	}
	return nil
}
