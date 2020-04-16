package version

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/90poe/connectctl/internal/ctl"
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager"
)

type versionCmdParams struct {
	ClusterURL string
}

// Command creates the the management commands
func Command() *cobra.Command {
	params := &versionCmdParams{}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return doVersion(cmd, params)
		},
	}

	ctl.AddClusterFlag(versionCmd, false, &params.ClusterURL)

	return versionCmd
}

func doVersion(cmd *cobra.Command, params *versionCmdParams) error {
	var (
		clusterInfo *connect.ClusterInfo
		err         error
	)

	if params.ClusterURL != "" {
		clusterInfo, err = getClusterInfo(params.ClusterURL)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Commit: %s\n", version.GitHash)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
	fmt.Printf("GO Version: %s\n", runtime.Version())
	fmt.Printf("GOOS: %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)

	if clusterInfo != nil {
		fmt.Printf("Connect Worker Version: %s\n", clusterInfo.Version)
	}

	return nil
}

func getClusterInfo(clusterURL string) (*connect.ClusterInfo, error) {
	config := &manager.Config{
		ClusterURL: clusterURL,
		Version:    version.Version,
	}

	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(config.ClusterURL, connect.WithUserAgent(userAgent))
	if err != nil {
		return nil, errors.Wrap(err, "error creating connect client")
	}

	mngr, err := manager.NewConnectorsManager(client, config)
	if err != nil {
		return nil, errors.Wrap(err, "error creating connectors manager")
	}

	clusterInfo, err := mngr.GetClusterInfo()
	if err != nil {
		return nil, errors.Wrap(err, "error getting cluster info")
	}

	return clusterInfo, nil
}
