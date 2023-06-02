package collect

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"time"
)

type ContributorUpload struct {
	RepoInfo struct {
		SubtaskId string `json:"subtask_id"`
	} `json:"repo_info"`
	LastCommitter struct {
		CommitterName  string    `json:"committer_name"`
		CommitterEmail string    `json:"committer_email"`
		CommitDate     time.Time `json:"commit_date"`
		CommitMessage  string    `json:"commit_message"`
	} `json:"last_committer"`
	Committers []Contributor `json:"committers"`
}

type Contributor struct {
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	LastCommitDate time.Time `json:"last_commit_date"`
	CommitCount    int       `json:"commit_count"`
}

func CollectDir(ctx context.Context, dir string) (*ContributorUpload, error) {
	logger := logctx.Use(ctx).Sugar().Named("collector")
	logger.Debugf("plain open %s", dir)
	repo, e := git.PlainOpen(dir)
	if errors.Is(e, git.ErrRepositoryNotExists) {
		logger.Debugf("no repository found")
		return nil, git.ErrRepositoryNotExists
	}
	if e != nil {
		logger.Debugf("open repository failed: %s", e.Error())
		return nil, fmt.Errorf("collector open repository: %w", e)
	}
	head, e := repo.Head()
	if e != nil {
		logger.Debugf("get head failed: %s", e.Error())
		return nil, fmt.Errorf("collector get head failed: %w", e)
	}
	var r = &ContributorUpload{
		Committers: nil,
	}
	set := make(map[[2]string]*Contributor)

	commit, e := repo.CommitObject(head.Hash())
	if e != nil {
		return nil, e
	}
	r.LastCommitter.CommitterEmail = commit.Author.Email
	r.LastCommitter.CommitterName = commit.Author.Name
	r.LastCommitter.CommitDate = commit.Author.When
	r.LastCommitter.CommitMessage = commit.Message
	for counter := 0; counter < 2000; counter++ {
		if errors.Is(e, plumbing.ErrObjectNotFound) || errors.Is(e, object.ErrParentNotFound) {
			break
		}
		if e != nil {
			logger.Warnf("errors during iterate commits, %s", e.Error())
			break
		}
		key := [2]string{commit.Author.Name, commit.Author.Email}
		if _, ok := set[key]; !ok {
			contributor := &Contributor{
				CommitterName:  commit.Author.Name,
				CommitterEmail: commit.Author.Email,
				LastCommitDate: commit.Author.When,
				CommitCount:    0,
			}
			set[key] = contributor
		}
		set[key].CommitCount++
		commit, e = commit.Parent(0)
	}

	for _, contributor := range set {
		r.Committers = append(r.Committers, *contributor)
	}
	return r, nil
}
