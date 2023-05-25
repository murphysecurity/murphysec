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
	LastCommitter *Contributor  `json:"last_committer"`
	Committers    []Contributor `json:"committers"`
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
		return nil, nil
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

	set := make(map[[2]string]*Contributor)
	var last *Contributor

	commit, e := repo.CommitObject(head.Hash())
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
			if last == nil {
				last = contributor
			}
		}
		set[key].CommitCount++
		commit, e = commit.Parent(0)
	}

	var r = &ContributorUpload{
		LastCommitter: last,
		Committers:    nil,
	}
	for _, contributor := range set {
		r.Committers = append(r.Committers, *contributor)
	}
	return r, nil
}
