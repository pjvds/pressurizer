package main

import (
	g "github.com/libgit2/git2go"
	"testing"
)

func TestCanOpenRepository(t *testing.T) {
	_, err := g.OpenRepository(".")

	if err != nil {
		t.Errorf("unable to open repository: %v", err)
	}
}

func TestCanGetHead(t *testing.T) {
	repo, _ := g.OpenRepository(".")
	ref, _ := repo.LookupReference("HEAD")

	if ref.Name() != "HEAD" {
		t.Errorf("expected ref name to be HEAD instead of %s", ref.Name())
	}
}
