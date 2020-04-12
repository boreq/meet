package commands

import (
	"github.com/boreq/guinea"
	"github.com/boreq/meet/cmd/meet/commands/users"
)

var MainCmd = guinea.Command{
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"run":    &runCmd,
		"config": &defaultConfigCmd,
		"users":  &users.UsersCmd,
	},
	ShortDescription: "a hydroponic farming infrastructure management service",
	Description: `
Hydro manages your hydroponic farming infrastructure.
`,
}

func runMain(c guinea.Context) error {
	return guinea.ErrInvalidParms
}
