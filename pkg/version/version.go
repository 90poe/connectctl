package version

import "fmt"

// Version specifies the version
var Version string

// BuildDate specifies the date
var BuildDate string

// GitTag is the git tag used when building
var GitTag string

// GitHash specifies the cimmit has associated with the build
var GitHash string

// ToString will convert the version information to a string
func ToString() string {
	return fmt.Sprintf("Version: %s, Build Date: %s, Git Tag: %s, Git Commit: %s", Version, BuildDate, GitTag, GitHash)
}
