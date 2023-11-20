package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/jessevdk/go-flags"

	"github.com/netsec-ethz/bootstrapper"
	"github.com/netsys-lab/scion-host/environment"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	Bootstrap string `short:"b" long:"bootstrap" description:"Bootstrapping URL" optional:"yes"`
}

var opts Options

func main() {

	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	endhostEnv := environment.EndhostEnv
	fmt.Println("Got env")
	fmt.Println(endhostEnv)
	endhostEnv.BootstrappingUrl = opts.Bootstrap

	// TODO: Proper error handling, do not fatal in here...
	endhostEnv.Install()

	code := bootstrapper.Run(filepath.Join(endhostEnv.ConfigPath, "bootstrapper.toml"), endhostEnv.ConfigPath)
	if code != 0 {
		log.Fatal("Bootstrapping failed")
	}

	// TODO: Supervise processes, ensure everything is running, restart in case something crashes
	// TODO: Write to logs and give helpful output
	dispatcherCmd := exec.Command(endhostEnv.DispatcherBinaryPath, "--config", filepath.Join(endhostEnv.DispatcherConfigPath, "dispatcher.toml"))
	dispatcherCmd.Stderr = os.Stderr
	dispatcherCmd.Stdout = os.Stdout

	daemonCmd := exec.Command(endhostEnv.DaemonBinaryPath, "--config", filepath.Join(endhostEnv.DaemonConfigPath, "sciond.toml"))
	daemonCmd.Stderr = os.Stderr
	daemonCmd.Stdout = os.Stdout

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := dispatcherCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = dispatcherCmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := daemonCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = daemonCmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
