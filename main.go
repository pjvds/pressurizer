package main

import (
	"flag"
	"fmt"
	"github.com/libgit2/git2go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var path string = *flag.String("path", ".", "the path of the directory that contains the git repository")

func main() {
	flag.Parse()

	cancel := trapSignal()

	run(cancel)
}

func run(cancel <-chan (bool)) {
	stop := false
	repo, err := git.OpenRepository(path)

	if err != nil {
		panic(err.Error())
	}

	for !stop {
		head, err := repo.LookupReference("HEAD")
		ref, err := head.Resolve()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(ref.Target())
		}

		stop = wait(cancel, 1*time.Second)
	}
}

func wait(cancel <-chan (bool), duration time.Duration) bool {
	cancelled := false

	select {
	case <-cancel:
		cancelled = true
	case <-time.After(duration):
	}

	return cancelled
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
