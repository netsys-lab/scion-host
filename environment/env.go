package environment

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/netsys-lab/scion-host/pkg/fileops"
)

var EndhostEnv *EndhostEnvironment

type EndhostEnvironment struct {
	BasePath             string
	ConfigPath           string
	DispatcherBinaryPath string
	DispatcherConfigPath string
	DaemonBinaryPath     string
	DaemonConfigPath     string
	EmbeddedFiles        fs.FS
	EmbeddedFolder       string
}

func (endhostEnv *EndhostEnvironment) Install() {
	if _, err := os.Stat(endhostEnv.BasePath); err != nil {
		err := os.MkdirAll(endhostEnv.BasePath, 0777)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Copying folder")
		err = fileops.CopyDir(endhostEnv.BasePath, endhostEnv.EmbeddedFiles, endhostEnv.EmbeddedFolder)
		if err != nil {
			log.Fatal(err)
		}

		err = fileops.ReplaceStringInFile(filepath.Join(endhostEnv.ConfigPath, "sciond.toml"), "{configDir}", endhostEnv.ConfigPath)
		if err != nil {
			log.Fatal("Failed to configure folder in sciond.toml: ", err)
		}

		err = fileops.ReplaceStringInFile(filepath.Join(endhostEnv.ConfigPath, "dispatcher.toml"), "{configDir}", endhostEnv.ConfigPath)
		if err != nil {
			log.Fatal("Failed to configure folder in sciond.toml: ", err)
		}

		// TODO: This could also be only windows specific
		err = fileops.ReplaceSingleBackslashWithDouble(filepath.Join(endhostEnv.ConfigPath, "dispatcher.toml"))
		if err != nil {
			log.Fatal(err)
		}
		err = fileops.ReplaceSingleBackslashWithDouble(filepath.Join(endhostEnv.ConfigPath, "sciond.toml"))
		if err != nil {
			log.Fatal(err)
		}
		// {configDir} ->
	}
}
