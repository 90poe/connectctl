package version

import (
	"runtime"

	log "github.com/sirupsen/logrus"
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
	log.Infof("Version: %s\n", version.Version)
	log.Infof("Commit: %s\n", version.GitHash)
	log.Infof("Build Date: %s\n", version.BuildDate)
	log.Infof("GO Version: %s\n", runtime.Version())
	log.Infof("GOOS: %s\n", runtime.GOOS)
	log.Infof("GOARCH: %s\n", runtime.GOARCH)
}
