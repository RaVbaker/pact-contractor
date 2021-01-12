package speccontext

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
)

const (
	GitFlowDevelopBranch = "develop"
	legacyMasterName     = "master"
	gitPath              = "./.git"
)

type GitContext struct {
	Context
	Author    string
	CommitSHA string
	Branch    string
}

func NewGitContext(specTag, origin, author, branch, commitSHA string) GitContext {
	gitContext := extractGitContext(commitSHA)

	// overwrite extracted values
	if len(author) != 0 {
		gitContext.Author = author
	}
	if len(branch) != 0 {
		gitContext.Branch = normalizeBranchName(branch)
	}
	if len(commitSHA) > 0 && len(gitContext.CommitSHA) == 0 {
		gitContext.CommitSHA = commitSHA
	}
	if specTag == BranchSpecTag {
		if len(gitContext.Branch) != 0 {
			specTag = gitContext.Branch
		} else {
			specTag = DefaultSpecTag
		}
	}

	gitContext.Context = NewContext(specTag, origin)
	return gitContext
}

func CurrentBranchName() string {
	gitContext := extractGitContext("")
	return gitContext.Branch
}

func extractGitContext(commitSHA string) GitContext {
	if ok, err := afero.DirExists(fs, gitPath); !ok || err != nil {
		log.Printf("No GIT repository found under %s", gitPath)
		return GitContext{}
	}

	commit, branchName := fetchGitDetails(gitPath, commitSHA)
	if commit == nil {
		log.Printf("No commits found for HEAD/%s under %s", commitSHA, gitPath)
		return GitContext{}
	}
	return GitContext{Branch: branchName, CommitSHA: commit.Hash.String(), Author: commit.Author.Name}
}

func fetchGitDetails(gitPath, commitSHA string) (*object.Commit, string) {
	r, err := git.PlainOpen(gitPath)
	if err != nil {
		log.Printf("Git repository open error: %v", err)
		return nil, ""
	}

	ref, err := findReferenceFromCommitSHA(r, commitSHA)
	if err != nil {
		log.Printf("Git reference resolve[%s] error: %v", commitSHA, err)
		return nil, ""
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Printf("Git commit[%v] open error: %v", ref.Hash(), err)
		return nil, ""
	}

	branchName := normalizeBranchName(ref.Name().Short())
	return commit, branchName
}

func findReferenceFromCommitSHA(r *git.Repository, commitSHA string) (*plumbing.Reference, error) {
	if len(commitSHA) == 0 {
		return r.Head()
	}

	hash, err := r.ResolveRevision(plumbing.Revision(commitSHA))
	if err != nil {
		return nil, err
	}

	refs, _ := r.References()
	var ref *plumbing.Reference
	err = refs.ForEach(func(iterRef *plumbing.Reference) error {
		if iterRef != nil && iterRef.Hash() == *hash {
			ref = iterRef
			refs.Close()
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("reference %s not found when iterating, %w", commitSHA, err)
	}

	if ref == nil {
		log.Printf("Couldn't find reference for %q, expanded hash to: %q, use --git-branch to overwrite branch", commitSHA, hash.String())
		ref = plumbing.NewReferenceFromStrings("-branch-not-found-", hash.String())
	}

	return ref, err
}

func normalizeBranchName(branch string) string {
	if branch == legacyMasterName {
		branch = DefaultSpecTag
	}
	return branch
}
