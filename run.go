package main

import (
	"context"
	"flag"
	"io/ioutil"
	"syscall"

	"github.com/purposed/tidy/tidy"

	"github.com/purposed/good/process"
	"github.com/purposed/good/serialization"
	"github.com/sirupsen/logrus"
)

const (
	runCommandName = "run"
	runCommandArgs = "[-cfg config.json]"
	runCommandHelp = "Run the tidy watcher"
)

type runCommand struct {
	configPath string
}

func (cmd *runCommand) Name() string      { return runCommandName }
func (cmd *runCommand) Args() string      { return runCommandArgs }
func (cmd *runCommand) ShortHelp() string { return runCommandHelp }
func (cmd *runCommand) LongHelp() string  { return runCommandHelp }
func (cmd *runCommand) Hidden() bool      { return false }

func (cmd *runCommand) Register(fs *flag.FlagSet) {
	fs.StringVar(&cmd.configPath, "cfg", "config.json", "The watcher configuration file.")
}

func (cmd *runCommand) Run(ctx context.Context, args []string) error {
	raw, err := ioutil.ReadFile(cmd.configPath)
	if err != nil {
		return err
	}

	var cfg tidy.Config
	if err := serialization.Unmarshal(raw, &cfg, serialization.JSON); err != nil {
		return err
	}

	engine, err := tidy.NewEngine(&cfg)
	if err != nil {
		return err
	}

	engine.Start()
	process.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)

	logrus.Info("asking engine to quit")
	engine.Stop()

	logrus.Info("goodbye")

	return nil
}
