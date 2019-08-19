//nolint:gochecknoglobals
package version

import "fmt"

// Version specifies the version
var Version string

// BuildDate specifies the date
var BuildDate string

// GitHash specifies the cimmit has associated with the build
var GitHash string

// ToString will convert the version information to a string
func ToString() string {
	return fmt.Sprintf("Version: %s, Build Date: %s, Git Commit: %s", Version, BuildDate, GitHash)
}
