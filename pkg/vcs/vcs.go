package vcs

type VCS interface {
	WithRepoPath(path string) (Repo, error)
	// CurrentRepo trying to find repo in current dir or parent directories
	CurrentRepo() (Repo, error)
}

type Repo interface {
	GetCurrentBranch() (string, error)
	GetRemotes() (map[string]string, error)
	HooksPath() string
}
