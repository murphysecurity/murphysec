package model

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/murphysecurity/murphysec/errors"
	giturls "github.com/whilp/git-urls"
	"time"
)

type GitInfo struct {
	RemoteName     string    `json:"remote_name"`
	RemoteURL      string    `json:"remote_url"`
	HeadCommitHash string    `json:"head_commit_hash"`
	HeadRefName    string    `json:"head_ref_name"`
	ProjectName    string    `json:"project_name"`
	CommitMsg      string    `json:"commit_msg"`
	Committer      string    `json:"committer"`
	CommitterEmail string    `json:"committer_email"`
	CommitTime     time.Time `json:"commit_time"`
}

func getGitInfo(dir string) (*GitInfo, error) {
	repo, e := git.PlainOpen(dir)
	if errors.Is(e, git.ErrRepositoryNotExists) ||
		errors.Is(e, git.ErrRepositoryIncomplete) {
		return nil, ErrNoGitRepo
	}
	if e != nil {
		return nil, fmt.Errorf("open git repo failed: %w", e)
	}
	// get remote
	remotes, e := repo.Remotes()
	if e != nil {
		return nil, fmt.Errorf("enumeration git remotes failed: %w", e)
	}
	var selectedRemote *git.Remote
	if len(remotes) == 0 {
		return nil, ErrNoGitRemoteFound
	}
	for _, it := range remotes {
		if it.Config().Name == "origin" {
			selectedRemote = it
			break
		}
	}
	if selectedRemote == nil {
		selectedRemote = remotes[0]
	}
	remoteUrls := selectedRemote.Config().URLs
	gitInfo := &GitInfo{
		RemoteName:     selectedRemote.Config().Name,
		RemoteURL:      "",
		HeadCommitHash: "",
		HeadRefName:    "",
		ProjectName:    "",
	}
	for _, it := range remoteUrls {
		u, e := giturls.Parse(it)
		if e != nil {
			continue
		}
		u.User = nil
		gitInfo.RemoteURL = u.String()
		gitInfo.ProjectName = u.Path
	}
	head, e := repo.Head()
	if e == nil && head != nil {
		commit, e := repo.CommitObject(head.Hash())
		if e == nil {
			gitInfo.CommitTime = commit.Committer.When
			gitInfo.CommitMsg = commit.Message
			gitInfo.Committer = commit.Committer.Name
			gitInfo.CommitterEmail = commit.Committer.Email
		}
		gitInfo.HeadCommitHash = head.Hash().String()
		gitInfo.HeadRefName = head.Name().String()
	}
	return gitInfo, nil
}

