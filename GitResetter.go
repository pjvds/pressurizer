package main

import (
	"github.com/libgit2/git2go"
)

func ResetRepository(path string) error {
	repo, err := git.OpenRepository(path)
	defer repo.Free()
	if err != nil {
		return err
	}

	headRef, err := repo.LookupReference("HEAD")
	defer headRef.Free()
	if err != nil {
		return err
	}

	ref, err := headRef.Resolve()
	defer ref.Free()
	if err != nil {
		return err
	}

	commit, err := repo.LookupCommit(ref.Target())
	defer commit.Free()
	if err != nil {
		return err
	}

	return repo.Checkout(&git.CheckoutOpts{
		Strategy: git.CHECKOUT_FORCE,
	})
}
