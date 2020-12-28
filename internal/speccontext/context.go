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
	branch, commitSHA := retrieveGitContextDetails()
	
	return GitContext{Context: Context{SpecTag: specTag}, CommitSHA: commitSHA, Branch: branch}
}

func retrieveGitContextDetails() (string, string) {
	r, err := git.PlainOpen("../.git")
	if err != nil {
		panic(err.Error())
	}
	
	rev, _ := r.Head()
	branch := strings.Replace(rev.Name().String(), "refs/heads/", "", 1)
	commitSHA := rev.Hash()
	return branch, commitSHA.String()
}