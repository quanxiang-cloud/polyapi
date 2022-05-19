package version

import (
	"fmt"
)

// polyapi version define
const (
	Version     = "v1.1.0"
	ReleaseDate = "2022-04-28"
	GitCommit   = "e0aeeec7"
)

// FullVersion get full version of polyapi
func FullVersion() string {
	return fmt.Sprintf("%s(%s@%s)", Version, ReleaseDate, GitCommit)
}
