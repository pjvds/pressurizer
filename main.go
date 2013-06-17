package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var repositoryPath string = *flag.String("path", ".", "the path of the directory that contains the git repository")

func main() {
	flag.Parse()

	cancel := trapSignal()
	watcher, err := NewGitWatcher(repositoryPath)
	defer watcher.Close()

	if err != nil {
		panic(err)
	}

	cancelled := false
	for !cancelled {
		select {
		case newHeadCommit := <-watcher.HeadChanged:
			fmt.Println("new head: ", newHeadCommit.Id())

		case err := <-watcher.Error:
			fmt.Println("error: ", err)

		case cancelled = <-cancel:
			fmt.Println("cancelled... will exit soon")

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
