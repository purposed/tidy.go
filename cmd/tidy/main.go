package main

import (
	"flag"

	"github.com/genuinetools/pkg/cli"
	"github.com/purposed/tidy/version"
)

func main() {
	p := cli.NewProgram()
	p.Name = "tidy"
	p.Description = "Automated filesystem organizer"
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	p.Commands = []cli.Command{
		&runCommand{},
	}

	p.FlagSet = flag.NewFlagSet("cli", flag.ExitOnError)

	p.Run()
}
