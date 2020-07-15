package config

type MirrorRepoConfig struct {
	ID           *int64
	FullName     *string
	Tracked      bool
	ShouldMirror bool
	Mirrored     bool
}

type MirrorConfig struct {
	GitHubOrg string
	GitLabOrg string
	Repos     []*MirrorRepoConfig
}
