package speccontext

import (
	"strings"
	
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Context struct {
	SpecTag string
	VerificationStatus bool
}

type GitContext struct {
	Context
	Author string
	CommitSHA string
	Branch string
}


func NewGitContext(specTag, author, branch, commitSHA string) GitContext {
	gitContext := extractGitContext(commitSHA)
	// overwrite extracted values
	if len(author) != 0 {
		gitContext.Author = author
	}
	if len(branch) != 0 {
		gitContext.Branch = branch
	}
	if specTag == BranchSpecTag {
		if len(gitContext.Branch) != 0 {
			specTag = gitContext.Branch
		} else {
			specTag = DefaultSpecTag
		}
	}
	
	gitContext.Context = NewContext(specTag)
	return gitContext
}

func NewContext(specTag string) Context {
	return Context{SpecTag: specTag}
}

func extractGitContext(commitSHA string) GitContext {
	commit, branchName := fetchGitDetails(commitSHA)
	return GitContext{Branch: branchName, CommitSHA: commit.Hash.String(), Author: commit.Author.Name}
}

func fetchGitDetails(commitSHA string) (*object.Commit, string) {
	r, err := git.PlainOpen(gitPath)
	if err != nil {
		panic(err.Error())
	}
	
	var ref *plumbing.Reference
	if len(commitSHA) == 0 {
		ref, err = r.Head()
	} else {
		ref, err = r.Reference(plumbing.ReferenceName(commitSHA), true)
	}
	if err != nil {
		panic(err.Error())
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		panic(err.Error())
	}
	
	branchName := extractBranchName(ref.Name().String())
	return commit, branchName
}

const (
	legacyMasterName = "master"
	defaultBranch = "main"
	BranchSpecTag = "{branch}"
	DefaultSpecTag  = "main"
    gitPath = "./.git"
)

func extractBranchName(refName string) (branch string) {
	branch = strings.Replace(refName, "refs/heads/", "", 1)
	if branch == legacyMasterName {
		branch = defaultBranch
	}
	return
}