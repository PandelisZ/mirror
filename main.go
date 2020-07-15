package main

import (
	"log"
	"os"

	"github.com/pandelisz/mirror/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "sync",
		Usage:  "sync a GitHub org and mirror to GitLab",
		Action: commands.Sync.Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "gh-token",
				Value: "",
				Usage: "GitHub token to use if not defined in the GH_TOKEN environment variable",
			},
			&cli.StringFlag{
				Name:     "gh-org",
				Value:    "",
				Usage:    "GitHub org to mirror",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
