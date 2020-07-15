package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/manifoldco/promptui"
	"github.com/tj/go-spin"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type command struct {
}

// Sync is the exported struct to be used in main
var Sync = command{}

func (a command) Action(c *cli.Context) error {

	GitHubToken := os.Getenv("GH_TOKEN")

	if GitHubToken == "" {
		if c.String("gh-token") == "" {
			return cli.Exit("GitHub token has not been set:\nGH_TOKEN environment variable not found", 1)
		}
		GitHubToken = c.String("gh-token")
	}

	SelectedGhOrg := c.String("gh-org")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	s := spin.New()
	quit := make(chan bool)
	spinner := func() {
		for i := 0; i < 100000; i++ {
			select {
			case <-quit:
				return
			default:
				fmt.Printf("\r  \033[36mRetrieving repositories\033[m %s ", s.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	go spinner()

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, SelectedGhOrg, opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	quit <- true

	prompt := promptui.Select{
		Label: "Select repositories to sync",
		Items: allRepos,
		Size:  30, //TODO: Make this dynamic based on terminal size
		Templates: &promptui.SelectTemplates{
			Label:    "{{.FullName}}",
			Selected: "{{.FullName}}    ✅",
			Active:   "▶️    {{.FullName}}",
			Inactive: "{{.FullName}}",
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return cli.Exit(fmt.Sprintf("Prompt failed %v\n", err), 1)
	}

	fmt.Printf("You choose %q\n", result)

	return nil
}
