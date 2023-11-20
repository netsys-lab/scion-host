package environment

import "embed"

//go:embed darwinarm64
var embeddedDarwinArmFiles embed.FS

func init() {
	EndhostEnv = &EndhostEnvironment{
		BasePath:             "/Applications/scion-host/",
		ConfigPath:           "/Applications/scion-host/darwinarm64/",
		DispatcherBinaryPath: "/Applications/scion-host/darwinarm64/dispatcher",
		DispatcherConfigPath: "/Applications/scion-host/darwinarm64/",
		DaemonBinaryPath:     "/Applications/scion-host/darwinarm64/daemon",
		DaemonConfigPath:     "/Applications/scion-host/darwinarm64/",
		EmbeddedFiles:        embeddedDarwinArmFiles,
		EmbeddedFolder:       "darwinarm64",
	}
}
