package commands

import (
	"encoding/json"
	"os"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/boreq/meet/internal/config"
	"github.com/boreq/meet/internal/wire"
)

var runCmd = guinea.Command{
	Run: runRun,
	Arguments: []guinea.Argument{
		{
			Name:        "config",
			Optional:    false,
			Multiple:    false,
			Description: "Path to the config file",
		},
	},
	ShortDescription: "starts a server",
}

func runRun(c guinea.Context) error {
	conf, err := loadConfig(c.Arguments[0])
	if err != nil {
		return errors.Wrap(err, "could not load the config")
	}

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	if err := service.Start(); err != nil {
		return errors.Wrap(err, "could not start a service")
	}

	return service.Wait()
}

func loadConfig(path string) (*config.Config, error) {
	conf := config.Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open the config file")
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(conf); err != nil {
		return nil, errors.Wrap(err, "json decoding failed")
	}

	return conf, nil
}
