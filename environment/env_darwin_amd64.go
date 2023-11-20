package environment

import "embed"

//go:embed darwinx64
var embeddedDarwinFiles embed.FS

func init() {
	EndhostEnv = &EndhostEnvironment{
		BasePath:             "/Applications/scion-host/",
		ConfigPath:           "/Applications/scion-host/darwinx64/",
		DispatcherBinaryPath: "/Applications/scion-host/darwinx64/dispatcher",
		DispatcherConfigPath: "/Applications/scion-host/darwinx64/",
		DaemonBinaryPath:     "/Applications/scion-host/darwinx64/daemon",
		DaemonConfigPath:     "/Applications/scion-host/darwinx64/",
		EmbeddedFiles:        embeddedDarwinFiles,
		EmbeddedFolder:       "darwinx64",
	}
}
