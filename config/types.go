package config

type MirrorRepoConfig struct {
	ID           *int64
	FullName     *string
	Name         *string
	Private      bool
	Tracked      bool
	ShouldMirror bool
	Mirrored     bool
	Description  *string
	URL          *string
}

type MirrorConfig struct {
	GitHubOrg   string
	GitLabGroup string
	Repos       []*MirrorRepoConfig
}
