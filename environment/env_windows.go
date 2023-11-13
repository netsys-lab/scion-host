//go:build windows
// +build windows

package environment

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed windowsx64
var embeddedWindowsFiles embed.FS

func init() {

	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		fmt.Println("The Program Files directory could not be found.")
		os.Exit(1)
	}

	basePath := filepath.Join(programFiles, "scion-host")
	configPath := filepath.Join(basePath, "windowsx64")

	EndhostEnv = &EndhostEnvironment{

		ConfigPath:           configPath,
		BasePath:             basePath,
		DispatcherBinaryPath: filepath.Join(configPath, "dispatcher.exe"),
		DispatcherConfigPath: configPath,
		DaemonBinaryPath:     filepath.Join(configPath, "daemon.exe"),
		DaemonConfigPath:     configPath,
		EmbeddedFiles:        embeddedWindowsFiles,
		EmbeddedFolder:       "windowsx64",
	}
}
