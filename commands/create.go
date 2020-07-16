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
	"github.com/pandelisz/mirror/config"
	"github.com/tj/go-spin"
	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"google.golang.org/api/sourcerepo/v1"
)

type create struct {
}

// Create is the exported struct to be used in main
var Create = create{}

func (a create) Debug(c *cli.Context) error {
	ctx := context.Background()
	sourcerepoService, err := sourcerepo.NewService(ctx)
	if err != nil {
		return err
	}

	res, err := sourcerepoService.Projects.Repos.List("projects/plumbus").Do()
	if err != nil {
		return err
	}
	for _, r := range res.Repos {
		fmt.Print(r)
	}

	return nil
}

// Google Cloud Source Repository
func gcsRepository(c *cli.Context, ghOrg string, config config.MirrorConfig) error {
	ctx := context.Background()
	sourcerepoService, err := sourcerepo.NewService(ctx)
	if err != nil {
		return err
	}

	projectName := "exampleproject"

	for i, repo := range config.Repos {

		p := &sourcerepo.Repo{
			Name: fmt.Sprintf("projects/%s/repos/github_%s_%s", projectName, ghOrg, *repo.Name),
			MirrorConfig: &sourcerepo.MirrorConfig{
				Url: *repo.URL,
			},
		}
		res, err := sourcerepoService.Projects.Repos.Create(fmt.Sprintf("projects/%s", projectName), p).Do()
		if err != nil {
			log.Printf("Failed to create project for %s :\n%v", *repo.FullName, err)
		} else {
			log.Printf("Created repository for %s", res.Name)
			config.Repos[i].Mirrored = true
		}
	}

	return nil
}

func (a create) Action(c *cli.Context) error {

	GitLabToken := os.Getenv("GL_TOKEN")

	if GitLabToken == "" {
		if c.String("gl-token") == "" {
			return cli.Exit("GitLab token has not been set:\nGL_TOKEN environment variable not found", 1)
		}
		GitLabToken = c.String("gl-token")
	}

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
	authenticatedGhUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return cli.Exit("Could not retrieve user from given GitHub token", 1)
	}

	var conf config.MirrorConfig
	configFile, err := os.Open(c.String("config"))
	if err != nil {
		return cli.Exit("Failed to load config file", 0)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &conf)

	gLab, err := gitlab.NewClient(GitLabToken)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create GitLab client %v", err), 1)
	}

	//TODO: sort out pagination
	nameSpaces, _, err := gLab.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{})
	var nameSpaceID int
	found := false
	for _, ns := range nameSpaces {
		if ns.Path == conf.GitLabGroup {
			nameSpaceID = ns.ID
			found = true
			break
		}
	}

	if !found {
		return cli.Exit(fmt.Sprintf("Could not find a group you belong to that matches %s", conf.GitLabGroup), 1)
	}

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
	for i, repo := range conf.Repos {
		p := &gitlab.CreateProjectOptions{
			Name:        repo.Name,
			Description: repo.Description,
			Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
			ImportURL:   gitlab.String(fmt.Sprintf("https://%s:%s@github.com/%s.git", *authenticatedGhUser.Login, GitHubToken, *repo.FullName)),
			NamespaceID: &nameSpaceID,
			Mirror:      gitlab.Bool(true),
		}
		_, _, err := gLab.Projects.CreateProject(p)
		if err != nil {
			log.Printf("Failed to create project for %s :\n%v", *repo.FullName, err)
		} else {
			log.Printf("Created repository for %s", *p.Name)
			conf.Repos[i].Mirrored = true
			conf.Repos[i].Replica = &config.ReplicaRepo{
				Name:       p.Name,
				CloneURL:   gitlab.String(fmt.Sprintf("https://gitlab.com/%s/%s.git", conf.GitLabGroup, *p.Name)),
				SSHURL:     gitlab.String(fmt.Sprintf("git@gitlab.com:%s/%s.git", conf.GitLabGroup, *p.Name)),
				ProjectURL: gitlab.String(fmt.Sprintf("https://gitlab.com/%s/%s", conf.GitLabGroup, *p.Name)),
			}
		}
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
