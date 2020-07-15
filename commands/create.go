package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pandelisz/mirror/config"
	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

type create struct {
}

// Create is the exported struct to be used in main
var Create = create{}

func (a create) Action(c *cli.Context) error {

	GitLabToken := os.Getenv("GL_TOKEN")

	if GitLabToken == "" {
		if c.String("gl-token") == "" {
			return cli.Exit("GitLab token has not been set:\nGL_TOKEN environment variable not found", 1)
		}
		GitLabToken = c.String("gl-token")
	}

	var config config.MirrorConfig
	configFile, err := os.Open(c.String("config"))
	if err != nil {
		return cli.Exit("Failed to load config file", 0)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &config)

	gLab, err := gitlab.NewClient(GitLabToken)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create GitLab client %v", err), 1)
	}

	//TODO: sort out pagination
	nameSpaces, _, err := gLab.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{})
	var nameSpaceID int
	found := false
	for _, ns := range nameSpaces {
		if ns.Path == config.GitLabGroup {
			nameSpaceID = ns.ID
			found = true
			break
		}
	}

	if !found {
		return cli.Exit(fmt.Sprintf("Could not find a group you belong to that matches %s", config.GitLabGroup), 1)
	}

	for _, repo := range config.Repos {
		p := &gitlab.CreateProjectOptions{
			Name:        repo.Name,
			Description: repo.Description,
			Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
			ImportURL:   repo.URL,
			NamespaceID: &nameSpaceID,
			Mirror:      gitlab.Bool(true),
		}
		_, _, err := gLab.Projects.CreateProject(p)
		if err != nil {
			log.Printf("Failed to create project for %s :\n%v", *repo.FullName, err)
		}
	}

	fmt.Printf("Successfully created repositories in GitLab!")
	return nil

}
