package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	git "github.com/libgit2/git2go/v30"
	"github.com/pandelisz/mirror/config"
	"github.com/tj/go-spin"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type sync struct {
}

// Sync is the exported struct to be used in main
var Sync = sync{}

func (a sync) Action(c *cli.Context) error {

	GitHubToken := os.Getenv("GH_TOKEN")

	if GitHubToken == "" {
		if c.String("gh-token") == "" {
			return cli.Exit("GitHub token has not been set:\nGH_TOKEN environment variable not found", 1)
		}
		GitHubToken = c.String("gh-token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// authenticatedGhUser, _, err := client.Users.Get(ctx, "")
	// if err != nil {
	// 	return cli.Exit("Could not retrieve user from given GitHub token", 1)
	// }

	var conf config.MirrorConfig
	configFile, err := os.Open(c.String("config"))
	if err != nil {
		return cli.Exit("Failed to load config file", 0)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &conf)

	s := spin.New()
	quit := make(chan bool)
	spinner := func() {
		for i := 0; i < 100000; i++ {
			select {
			case <-quit:
				return
			default:
				fmt.Printf("\r  \033[36mCreating repositories\033[m %s ", s.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	go spinner()

	for _, repo := range conf.Repos {

		if !repo.ShouldMirror || repo.Mirrored {
			continue
		}

		ghRepo, _, err := client.Repositories.Get(ctx, conf.GitHubOrg, *repo.Name)
		if err != nil {
			log.Printf("Failed to query project for %s :\n%v", *repo.FullName, err)
			continue
		}
		lastUpdatedUnix := ghRepo.UpdatedAt.Unix()

		// Don't bother updating if no changes since last sync
		if repo.Replica.LastSync > lastUpdatedUnix {
			log.Println("Skipped")
			continue
		}

		clonedRepository, err := git.Clone(*repo.Replica.SSHURL, "./clonedir/repo", &git.CloneOptions{
			Bare: true,
		})

		if err != nil {
			return err
		}

		fmt.Print(clonedRepository.Remotes)
	}
	quit <- true

	configFile.Close()
	configFileOut, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.String("out"), configFileOut, 0644)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully created repositories in GitLab!")
	return nil

}
