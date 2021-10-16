package vcs

import (
	"math"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"gohooks/pkg/common"
	"gohooks/pkg/executor"
)

const (
	gitFolder = ".git"
)

func NewGit(gitExecutor GitExecutor) VCS {
	return &git{executor: gitExecutor}
}

type git struct {
	executor GitExecutor
}

func (git *git) CurrentRepo() (Repo, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Finds out path to .git directory
	output, err := git.executor.Output([]string{"rev-parse", "--git-dir"}, executor.WithWorkdir(workDir))
	if err != nil {
		return nil, err
	}

	var repoPath string
	if outputStr := strings.Trim(string(output), " \n"); outputStr != gitFolder {
		// Receive path to repo without .git in path
		repoPath, _ = path.Split(string(output))
	} else {
		// We are at repo root
		repoPath = workDir
	}

	return &gitRepo{
		repoPath: repoPath,
		executor: git.executor,
	}, nil
}

func (git *git) WithRepoPath(repoPath string) (Repo, error) {
	exists, err := common.PathExists(path.Join(repoPath, gitFolder))
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.Errorf("git repo by repoPath %s not exists or it not git repo", repoPath)
	}

	return &gitRepo{
		repoPath: repoPath,
		executor: git.executor,
	}, nil
}

type gitRepo struct {
	repoPath string
	executor GitExecutor
}

func (repo *gitRepo) GetCurrentBranch() (string, error) {
	output, err := repo.output("rev-parse", "--abbrev-ref", "HEAD")
	// Remove unnecessary new line at the end
	return strings.Trim(string(output), " \n"), err
}

func (repo *gitRepo) GetRemotes() (map[string]string, error) {
	output, err := repo.output("remote", "-v")
	if err != nil {
		return nil, err
	}

	branchesToURLMap := map[string]string{}

	for i, s := range strings.Split(string(output), "\n") {
		if math.Mod(float64(i), 2) != 0 {
			continue
		}
		if s == "" {
			continue
		}
		// Fix unnecessary tabs in output of command
		s = regexp.MustCompile("\t+").ReplaceAllString(s, " ")

		parts := strings.Split(s, " ")
		if len(parts) != 3 {
			return nil, errors.New("unknown format for git remote -v")
		}

		branchesToURLMap[parts[0]] = parts[1]
	}
	return branchesToURLMap, nil
}

func (repo *gitRepo) HooksPath() string {
	return path.Join(repo.repoPath, gitFolder, "hooks")
}

//nolint:unused
func (repo *gitRepo) run(args ...string) error {
	return repo.executor.Run(args, executor.WithWorkdir(repo.repoPath))
}

func (repo *gitRepo) output(args ...string) ([]byte, error) {
	return repo.executor.Output(args, executor.WithWorkdir(repo.repoPath))
}
