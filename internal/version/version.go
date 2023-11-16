package version

import "runtime"

var (
	AppVersion string
	GitCommit  string
	GitBranch  string
	BuildDate  string
	GoVersion  = runtime.Version()
)
