package main

import (
	"github.com/libgit2/git2go"
)

func ResetRepository(path string) (commit *git.Commit, err error) {
	repo, err := git.OpenRepository(path)
	defer repo.Free()
	if err != nil {
		return
	}

	// headRef, err := repo.LookupReference("HEAD")
	// defer headRef.Free()
	// if err != nil {
	// 	return
	// }

	// ref, err := headRef.Resolve()
	// defer ref.Free()
	// if err != nil {
	// 	return
	// }

	// commit, err = repo.LookupCommit(ref.Target())
	// defer commit.Free()
	// if err != nil {
	// 	return
	// }

	// index, err := repo.Index()
	// defer index.Free()
	// if err != nil {
	// 	return
	// }

	// err = index.AddByPath(".")
	// if err != nil {
	// 	return
	// }

	// pressureRef, err := repo.CreateReference("refs/head/pressure", ref.Target(), true)
	// defer pressureRef.Free()
	// if err != nil {
	// 	return
	// }

	// treeId, err := index.WriteTree()
	// if err != nil {
	// 	return
	// }

	// signature := &git.Signature{
	// 	Name:  "Pressurizer",
	// 	Email: "pj@born2code.net",
	// 	When:  time.Now(),
	// }

	// tree, err := repo.LookupTree(treeId)
	// if err != nil {
	// 	return
	// }

	// commitId, err := repo.CreateCommit(pressureRef.Name(), signature, signature, "you where too late!", tree)
	// if err != nil {
	// 	return
	// }

	// commit, err = repo.LookupCommit(commitId)
	// if err != nil {
	// 	return
	// }

	err = repo.Checkout(&git.CheckoutOpts{
		Strategy: git.CHECKOUT_FORCE | git.CHECKOUT_REMOVE_UNTRACKED,
	})

	return
}
