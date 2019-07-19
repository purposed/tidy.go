package main

import (
	"flag"

	"github.com/genuinetools/pkg/cli"
	"github.com/purposed/fsclean/version"
)

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
