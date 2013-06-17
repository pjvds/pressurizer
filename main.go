package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/libgit2/git2go"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var repositoryPath string = *flag.String("path", ".", "the path of the directory that contains the git repository")

func main() {
	flag.Parse()

	cancel := trapSignal()
	repo, err := git.OpenRepository(repositoryPath)
	if err != nil {
		panic(err.Error())
	}

	headRef, _ := repo.LookupReference("HEAD")
	ref, _ := headRef.Resolve()
	fmt.Println(ref.Name())

	watcher, err := fsnotify.NewWatcher()

	watcher.Watch(path.Join(repo.Path(), "HEAD"))
	watcher.Watch(path.Join(repo.Path(), ref.Name()))

	cancelled := false
	for !cancelled {
		select {
		case event := <-watcher.Event:
			fmt.Println(event.String())
			break
		case error := <-watcher.Error:
			fmt.Println(error.Error())
		case cancelled = <-cancel:
			fmt.Println("cancelled... will exit soon")
			break
		case <-time.After(5 * time.Minute):
			fmt.Println("timed out!!!!!")
		}
	}
}

func trapSignal() <-chan (bool) {
	interrupt_chan := make(chan os.Signal, 2)
	cancelled := make(chan bool)

	signal.Notify(interrupt_chan, os.Interrupt)
	signal.Notify(interrupt_chan, syscall.SIGTERM)

	go func() {
		sig := <-interrupt_chan
		fmt.Printf("\nsignal %s received, cancelling...", sig)
		cancelled <- true
	}()

	return cancelled
}
