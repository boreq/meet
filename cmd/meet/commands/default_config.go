package commands

import (
	"encoding/json"
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/boreq/meet/internal/config"
)

var defaultConfigCmd = guinea.Command{
	Run:              runDefaultConfig,
	ShortDescription: "prints the default configuration",
}

func runDefaultConfig(c guinea.Context) error {
	conf := config.Default()

	j, err := json.MarshalIndent(conf, "", "   ")
	if err != nil {
		return errors.Wrap(err, "could not marshal the config")
	}

	fmt.Println(string(j))

	return nil
}
