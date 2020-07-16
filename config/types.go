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
	CloneURL     *string
	SSHURL       *string
	Replica      *ReplicaRepo `json:"replica,omitempty"`
}

type ReplicaRepo struct {
	Name       *string
	CloneURL   *string
	ProjectURL *string
	SSHURL     *string
	LastSync   int64 `json:"lastSync,omitempty"`
}

type MirrorConfig struct {
	GitHubOrg   string
	GitLabGroup string
	LastSync    int64 `json:"lastSync,omitempty"`
	Repos       []*MirrorRepoConfig
}
