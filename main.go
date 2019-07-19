package main

import (
	"flag"

	"github.com/dalloriam/fsclean/version"
	"github.com/genuinetools/pkg/cli"
)

const configFile = "config.json"

func main() {
	p := cli.NewProgram()
	p.Name = "fsclean"
	p.Description = "Configurable file cleaner"
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	p.Commands = []cli.Command{
		&runCommand{},
	}

	p.FlagSet = flag.NewFlagSet("cli", flag.ExitOnError)

	p.Run()
}
