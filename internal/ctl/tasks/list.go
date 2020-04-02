package tasks

import (
	"fmt"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

type ListOptions struct {
	// Use pointer here to make reference to the 'tasks' options.
	// Otherwise, the local copy of the config will be created an no
	// value from the flags will be set. It is caused by the global
	// scope nature of viper and multiple overrides of cluster flag
	// from the other commands.
	*GenericOptions

	Output string
}

func NewListCommand(options *GenericOptions) *cobra.Command {
	var opts = ListOptions{
		GenericOptions: options,
	}

	var cmd = cobra.Command{
		Use:   "list",
		Short: "Displays a list of tasks currently running for the connector.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return opts.Run(cmd.OutOrStdout())
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Output, "output", "o", "json", "The output format. Possible values are 'json' and 'table'")

	return &cmd
}

func (o *ListOptions) Run(out io.Writer) error {
	client, err := o.CreateClient(o.ClusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create http client")
	}
	tasks, _, err := client.GetConnectorTasks(o.ConnectorName)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve tasks")
	}
	return o.writeOutput(tasks, out)
}

func (o *ListOptions) writeOutput(tasks []connect.Task, out io.Writer) error {
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
	output, err := outputFn(tasks)
	if err != nil {
		return fmt.Errorf("failed to form output for '%s' type", outputType)
	}
	if _, err := out.Write(output); err != nil {
		return errors.Wrap(err, "failed to write output")
	}
	return nil
}
