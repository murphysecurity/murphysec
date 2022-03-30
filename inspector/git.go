package inspector

import (
	"container/list"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
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

func (g *GitInfo) ApiVo() *api.VoGitInfo {
	if g == nil {
		return nil
	}
	return &api.VoGitInfo{
		Commit:        g.HeadCommitHash,
		GitRef:        g.HeadRefName,
		GitRemoteUrl:  g.RemoteURL,
		CommitMessage: g.CommitMsg,
		CommitEmail:   g.CommitterEmail,
		CommitTime:    g.CommitTime,
	}
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
	logger.Debug.Printf("Selected remote: %s", selectedRemote.String())
	logger.Debug.Printf("Total %d urls", len(remoteUrls))
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
			logger.Debug.Printf("Parse git url failed: %s, url: %s", e.Error(), it)
			continue
		}
		u.User = nil
		gitInfo.RemoteURL = u.String()
		gitInfo.ProjectName = u.Path
	}
	head, e := repo.Head()
	if e != nil {
		logger.Warn.Println("Get HEAD failed.", e.Error())
	} else {
		if head != nil {
			commit, e := repo.CommitObject(head.Hash())
			if e == nil {
				gitInfo.CommitTime = commit.Committer.When
				gitInfo.CommitMsg = commit.Message
				gitInfo.Committer = commit.Committer.Name
				gitInfo.CommitterEmail = commit.Committer.Email
			} else {
				logger.Warn.Println("Get commit info failed.", e.Error())
			}
			gitInfo.HeadCommitHash = head.Hash().String()
			gitInfo.HeadRefName = head.Name().String()
		} else {
			logger.Warn.Println("HEAD is null")
		}
	}
	return gitInfo, nil
}
func CollectContributor(dir string) ([]api.Contributor, error) {
	repo, e := git.PlainOpen(dir)
	if e != nil {
		return nil, errors.Wrap(e, "open repo failed")
	}
	branches, e := repo.References()
	if e != nil {
		return nil, errors.Wrap(e, "list branches failed")
	}
	distinctHashSet := map[string]struct{}{}
	contributorSet := map[api.Contributor]struct{}{}
	e = branches.ForEach(func(reference *plumbing.Reference) error {
		h := reference.Hash()
		q := list.New()
		q.PushBack(h)
		for q.Len() > 0 {
			cur := q.Front().Value.(plumbing.Hash)
			q.Remove(q.Front())
			if len(distinctHashSet) > 10000 {
				return nil
			}
			if cur.IsZero() {
				continue
			}
			if _, ok := distinctHashSet[h.String()]; ok {
				continue
			}
			distinctHashSet[h.String()] = struct{}{}
			commit, e := repo.CommitObject(h)
			if e == plumbing.ErrObjectNotFound {
				continue
			}
			if e != nil {
				return errors.Wrap(e, "get commit failed")
			}
			for _, hash := range commit.ParentHashes {
				q.PushBack(hash)
			}
			if commit.Committer.Name != "" {
				contributorSet[api.Contributor{
					Name:  commit.Committer.Name,
					Email: commit.Committer.Email,
				}] = struct{}{}
			}
			if commit.Author.Name != "" {
				contributorSet[api.Contributor{
					Name:  commit.Author.Name,
					Email: commit.Author.Email,
				}] = struct{}{}
			}
		}
		return nil
	})
	if e != nil {
		return nil, errors.Wrap(e, "iterate failed.")
	}
	var rs []api.Contributor
	for contributor := range contributorSet {
		rs = append(rs, contributor)
	}
	return rs, e
}
