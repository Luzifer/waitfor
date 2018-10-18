package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig"
)

var (
	cfg = struct {
		CheckInterval  time.Duration `flag:"check-interval,i" default:"1s" description:"How long to wait after an unsuccessful check"`
		CommandTimeout time.Duration `flag:"command-timeout,c" default:"0" description:"Stop the command execution after this time"`
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		Shell          string        `flag:"shell,s" default:"/bin/bash" description:"Shell to execute with the given command (must accept -c flag)"`
		WaitTimeout    time.Duration `flag:"wait-timeout,w" default:"0" description:"Stop waiting for the command after this time"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("waitfor %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	started := time.Now()

	for {
		if err := executeCommand(rconfig.Args()[1:]); err == nil {
			break
		} else {
			log.WithField("error", err.Error()).Debug("Command was not successful")
		}

		if cfg.WaitTimeout > 0 && started.Add(cfg.WaitTimeout).Before(time.Now()) {
			log.WithField("timeout", cfg.WaitTimeout).Fatal("Wait timed out")
		}

		time.Sleep(cfg.CheckInterval)
	}

	log.Info("Command exited successful")
}

func executeCommand(cmdStr []string) error {
	cmd := exec.Command(cfg.Shell, "-c", strings.Join(cmdStr, " "))

	var (
		c                         = make(chan error, 1)
		ctx                       = context.Background()
		cancel context.CancelFunc = func() {}
	)
	if cfg.CommandTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), cfg.CommandTimeout)
	}

	go func() {
		c <- cmd.Run()
	}()

	for {
		select {
		case <-ctx.Done():
			cmd.Process.Kill()
			return errors.New("Command execution timed out")
		case err := <-c:
			cancel()
			return err
		}
	}
}
