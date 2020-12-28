package speccontext

import (
	"strings"
	
	"github.com/go-git/go-git/v5"
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


func NewGitContext(specTag string) GitContext {
	gitContext := lookupGitContext()
	gitContext.Context = Context{SpecTag: specTag}
	return gitContext
}

func lookupGitContext() GitContext {
	r, err := git.PlainOpen("../.git")
	if err != nil {
		panic(err.Error())
	}
	
	ref, _ := r.Head()
	commit, _ := r.CommitObject(ref.Hash())
	
	branchName := extractBranchName(ref.Name().String())
	return GitContext{Branch: branchName, CommitSHA: commit.Hash.String(), Author: commit.Author.Name}
}

const (
	legacyMasterName = "master"
	defaultBranch = "main"
)

func extractBranchName(refName string) (branch string) {
	branch = strings.Replace(refName, "refs/heads/", "", 1)
	if branch == legacyMasterName {
		branch = defaultBranch
	}
	return
}