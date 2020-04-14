package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/90poe/connectctl/internal/version"
)

// Command creates the the management commands
func Command() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "",
		Run:   doVersion,
	}

	return versionCmd
}

func doVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Commit: %s\n", version.GitHash)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
	fmt.Printf("GO Version: %s\n", runtime.Version())
	fmt.Printf("GOOS: %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
}
