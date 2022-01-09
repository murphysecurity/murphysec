package inspector

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"
	"murphysec-cli-simple/logger"
)

type GitInfo struct {
	RemoteName     string `json:"remote_name"`
	RemoteURL      string `json:"remote_url"`
	HeadCommitHash string `json:"head_commit_hash"`
	HeadRefName    string `json:"head_ref_name"`
}

func getGitInfo(dir string) (*GitInfo, error) {
	logger.Debug.Println("Try open git:", dir)
	repo, e := git.PlainOpen(dir)
	if e == git.ErrRepositoryNotExists {
		return nil, nil
	}
	if e == git.ErrRepositoryIncomplete {
		return nil, errors.New("Git repo incomplete, skip")
	}
	if e != nil {
		return nil, errors.Wrap(e, "Git err")
	}
	// get remote
	remotes, e := repo.Remotes()
	if e != nil {
		return nil, errors.Wrap(e, "Enumeration git remotes failed")
	}
	var selectedRemote *git.Remote
	logger.Debug.Println(fmt.Sprintf("Found %d remotes", len(remotes)))
	if len(remotes) == 0 {
		return nil, errors.New("No git remote found")
	}
	for _, it := range remotes {
		if it.Config().Name == "origin" {
			selectedRemote = it
			logger.Debug.Println(fmt.Sprintf("Remote: origin found"))
			break
		}
	}
	if selectedRemote == nil {
		selectedRemote = remotes[0]
		logger.Debug.Println("No origin remote, use first one")
	}
	remoteUrls := selectedRemote.Config().URLs
	logger.Debug.Println(fmt.Sprintf("Selected remote: %s", selectedRemote.String()))
	logger.Debug.Println(fmt.Sprintf("Total %d urls", len(remoteUrls)))
	gitInfo := &GitInfo{
		RemoteName:     selectedRemote.Config().Name,
		RemoteURL:      "",
		HeadCommitHash: "",
		HeadRefName:    "",
	}
	for _, it := range remoteUrls {
		u, e := giturls.Parse(it)
		if e != nil {
			logger.Debug.Println(fmt.Sprintf("Parse git url failed: %s, url: %s", e.Error(), it))
			continue
		}
		u.User = nil
		gitInfo.RemoteURL = u.String()
	}
	head, e := repo.Head()
	if e != nil {
		logger.Warn.Println("Get HEAD failed.", e.Error())
	} else {
		if head != nil {
			gitInfo.HeadCommitHash = head.Hash().String()
			gitInfo.HeadRefName = head.Name().String()
		} else {
			logger.Warn.Println("HEAD is null")
		}
	}
	return gitInfo, nil
}
