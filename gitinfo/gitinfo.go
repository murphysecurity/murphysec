package gitinfo

import (
	"context"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"time"
)

func GetSummary(ctx context.Context, dir string) (*Summary, error) {
	var (
		e      error
		logger = logctx.Use(ctx).Sugar()
	)

	repo, e := git.PlainOpen(dir)
	if errors.Is(e, git.ErrRepositoryNotExists) {
		return nil, ErrNoRepoFound
	}
	if e != nil {
		return nil, e
	}

	var summary Summary

	summary.RemoteAddr, e = getRemoteURL(ctx, repo)
	if e != nil {
		logger.Warnf("get remote url: %v", e)
	}
	head, e := repo.Head()
	if e != nil {
		logger.Warnf("get head: %v", e)
	} else {
		summary.BranchName = head.Name().Short()
		cinfo, e := getCommitInfo(ctx, repo, head.Hash())
		if e != nil {
			logger.Warnf("get commit: %v", e)
		} else {
			summary.CommitHash = cinfo.Hash
			summary.CommitMessage = cinfo.Message
			summary.CommitTime = cinfo.Time
			summary.AuthorEmail = cinfo.AuthorOrCommitterEmail
			summary.AuthorName = cinfo.AuthorOrCommitterName
		}
	}
	return &summary, nil
}

func getRemoteURL(ctx context.Context, repo *git.Repository) (string, error) {
	remotes, e := repo.Remotes()
	if e != nil {
		return "", e
	}

	var selectedRemote *git.Remote
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			selectedRemote = remote
			break // use origin if found
		}
		if selectedRemote == nil {
			selectedRemote = remote
		}
	}

	if selectedRemote == nil {
		return "", _ErrNoRemoteURLFound
	}

	var candidateURLs = selectedRemote.Config().URLs
	if len(candidateURLs) == 0 {
		return "", _ErrNoRemoteURLFound
	}

	return candidateURLs[0], nil
}

func getCommitInfo(ctx context.Context, repo *git.Repository, hash plumbing.Hash) (*commitInfo, error) {
	commit, e := repo.CommitObject(hash)
	if e != nil {
		return nil, e
	}
	var info commitInfo
	info.Message = commit.Message
	info.Hash = commit.Hash.String()
	info.AuthorOrCommitterEmail = commit.Author.Email
	info.AuthorOrCommitterName = commit.Author.Name
	info.Time = commit.Author.When
	if info.AuthorOrCommitterName == "" {
		info.AuthorOrCommitterEmail = commit.Committer.Email
		info.AuthorOrCommitterName = commit.Committer.Name
		info.Time = commit.Committer.When
	}
	return &info, nil
}

type commitInfo struct {
	Hash                   string
	AuthorOrCommitterName  string
	AuthorOrCommitterEmail string
	Time                   time.Time
	Message                string
}

type Summary struct {
	RemoteAddr    string
	BranchName    string
	CommitHash    string
	CommitMessage string
	CommitTime    time.Time
	AuthorName    string
	AuthorEmail   string
}
