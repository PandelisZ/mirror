package main

import (
	"log"
	"os"

	"github.com/pandelisz/mirror/commands"
	"github.com/urfave/cli/v2"
)

type Myflags struct {
	ghOrg   *cli.StringFlag
	glGroup *cli.StringFlag
	ghToken *cli.StringFlag
	glToken *cli.StringFlag
	out     *cli.StringFlag
	config  *cli.StringFlag
}

var flags Myflags = Myflags{
	ghOrg: &cli.StringFlag{
		Name:     "gh-org",
		Value:    "",
		Usage:    "GitHub org to mirror",
		Required: true,
	},
	glGroup: &cli.StringFlag{
		Name:     "gl-group",
		Value:    "",
		Usage:    "GitLab group to create mirrors inside of",
		Required: true,
	},
	ghToken: &cli.StringFlag{
		Name:  "gh-token",
		Value: "",
		Usage: "GitHub token to use if not defined in the GH_TOKEN environment variable",
	},
	glToken: &cli.StringFlag{
		Name:  "gl-token",
		Value: "",
		Usage: "GitLab token to use if not defined in the GL_TOKEN environment variable",
	},
	config: &cli.StringFlag{
		Name:      "config",
		Aliases:   []string{"f", "c", "file"},
		Value:     "repositories.json",
		TakesFile: true,
		Usage:     "Path to configuration file",
	},
	out: &cli.StringFlag{
		Name:    "out",
		Aliases: []string{"o"},
		Value:   "repositories.json",
		Usage:   "Name of repositories file",
	},
}

func main() {
	app := &cli.App{
		Name:  "mirror",
		Usage: "mirror a GitHub org to GitLab",
		Commands: []*cli.Command{
			{
				Name:   "import",
				Usage:  "Import your repos into a config file to be tracked",
				Action: commands.Import.Action,
				Flags: []cli.Flag{
					flags.ghToken,
					flags.ghOrg,
					flags.glGroup,
					flags.out,
				},
			},
			{
				Name:   "create",
				Usage:  "Create a config file with the repos you'd like to import",
				Action: commands.Create.Action,
				Flags: []cli.Flag{
					flags.ghToken,
					flags.glToken,
					flags.config,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
