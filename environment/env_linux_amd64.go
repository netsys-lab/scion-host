package environment

import "embed"

//go:embed linuxx64
var embeddedLinuxFiles embed.FS

func init() {
	EndhostEnv = &EndhostEnvironment{
		BasePath:             "/etc/scion-host/",
		ConfigPath:           "/etc/scion-host/linuxx64/",
		DispatcherBinaryPath: "/etc/scion-host/linuxx64/dispatcher",
		DispatcherConfigPath: "/etc/scion-host/linuxx64/",
		DaemonBinaryPath:     "/etc/scion-host/linuxx64/daemon",
		DaemonConfigPath:     "/etc/scion-host/linuxx64/",
		EmbeddedFiles:        embeddedLinuxFiles,
		EmbeddedFolder:       "linuxx64",
	}
}
