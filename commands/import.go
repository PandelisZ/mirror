package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/pandelisz/mirror/config"
	"github.com/tj/go-spin"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type command struct {
}

// Import is the exported struct to be used in main
var Import = command{}

func (a command) Action(c *cli.Context) error {

	GitHubToken := os.Getenv("GH_TOKEN")

	if GitHubToken == "" {
		if c.String("gh-token") == "" {
			return cli.Exit("GitHub token has not been set:\nGH_TOKEN environment variable not found", 1)
		}
		GitHubToken = c.String("gh-token")
	}

	SelectedGhOrg := c.String("gh-org")
	DesinationGlGroup := c.String("gl-group")

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

	var repoConfig []*config.MirrorRepoConfig
	for _, repo := range allRepos {
		repoConfig = append(repoConfig, &config.MirrorRepoConfig{
			ID:           repo.ID,
			Name:         repo.Name,
			Description:  repo.Description,
			FullName:     repo.FullName,
			URL:          repo.CloneURL,
			Tracked:      true,
			ShouldMirror: true,
			Mirrored:     false,
			Private:      true,
		})
	}

	mirrorConfig := config.MirrorConfig{
		GitHubOrg:   SelectedGhOrg,
		GitLabGroup: DesinationGlGroup,
		Repos:       repoConfig,
	}
	configFile, _ := json.MarshalIndent(mirrorConfig, "", "    ")
	_ = ioutil.WriteFile(c.String("out"), configFile, 0644)

	fmt.Printf("\nSuccessfully imported repositores for github.com/%s\n", SelectedGhOrg)
	return nil
}
