package main

import (
	"log"
	"os"

	"github.com/evenq/evenq-cli/src/eventCommand"
	"github.com/evenq/evenq-cli/src/imports"
	"github.com/evenq/evenq-cli/src/shared/api"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "login",
				Action: api.RunLogin,
			},
			{
				Name: "event",
				Subcommands: []*cli.Command{
					{
						Name:      "import",
						Usage:     "import data to an existing or new event",
						ArgsUsage: "CSV file, may be be gzipped",
						Action:    imports.Run,
					},
					{
						Name:   "create",
						Usage:  "create a new event",
						Action: eventCommand.RunCreate,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
