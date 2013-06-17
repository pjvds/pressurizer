package main

import (
	"github.com/howeyc/fsnotify"
	"github.com/libgit2/git2go"
)

type GitWatcher struct {
	HeadChanged chan *git.Commit // Channel for sending new head commits
	Error       chan *error

	fsWatcher *fsnotify.Watcher
	repo      *git.Repository // The repository to watch
	done      chan bool       // Channel for sending a quit message to the reader goroutine
	isClosed  bool            // Set to true when Close() is first called
}

func NewGitWatcher(path string) (*GitWatcher, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fsWatcher.Watch(repo.Path())

	gitWatcher := &GitWatcher{
		HeadChanged: make(chan *git.Commit),
		Error:       make(chan *error),
		fsWatcher:   fsWatcher,
		repo:        repo,
		done:        make(chan bool),
		isClosed:    false,
	}

	return gitWatcher, nil
}

func (w *GitWatcher) watchRepository() {
	defer w.repo.Free()
	defer w.fsWatcher.Close()

	previousCommitId, err := getHeadCommitId(w.repo)
	if err != nil {
		w.Error <- &err
	}

	for {
		select {
		case err := <-w.fsWatcher.Error:
			w.Error <- &err
		case <-w.fsWatcher.Event:
			commitId, err := getHeadCommitId(w.repo)
			if err != nil {
				w.Error <- &err
			}

			if previousCommitId.Cmp(commitId) != 0 {
				commit, err := w.repo.LookupCommit(commitId)
				if err != nil {
					w.Error <- &err
				}

				previousCommitId = commitId

				w.HeadChanged <- commit
				commit.Free()
			}
		case <-w.done:
			return
		}
	}
}

func getHeadCommitId(repo *git.Repository) (*git.Oid, error) {
	headRef, err := repo.LookupReference("HEAD")
	defer headRef.Free()
	if err != nil {
		return nil, err
	}

	ref, err := headRef.Resolve()
	defer ref.Free()
	if err != nil {
		return nil, err
	}

	return ref.Target(), nil
}

func (w *GitWatcher) Close() error {
	if w.isClosed {
		return nil
	}
	w.isClosed = true

	close(w.Error)
	close(w.HeadChanged)

	if w.fsWatcher != nil {
		w.fsWatcher.Close()
	}

	if w.repo != nil {
		w.repo.Free()
	}

	w.done <- true
	return nil
}
