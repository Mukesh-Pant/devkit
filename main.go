package main

import (
	"devkit/cmd"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Inject version information into cmd package
	cmd.SetVersionInfo(Version, Commit, BuildDate)
	
	cmd.Execute()
}
