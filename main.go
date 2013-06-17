package main

import (
	"flag"
	"fmt"
	"time"
)

var repositoryPath string = *flag.String("path", ".", "the path of the directory that contains the git repository")

func main() {
	flag.Parse()

	cancel := TrapInterrupts()
	watcher, err := NewGitWatcher(repositoryPath)
	defer watcher.Close()

	if err != nil {
		panic("unable to watch git repository at '" + repositoryPath + "': " + err.Error())
	}

	cancelled := false
	for !cancelled {
		select {
		case newHeadCommit := <-watcher.HeadChanged:
			fmt.Println("new head: ", newHeadCommit.Id())

		case err := <-watcher.Error:
			fmt.Println("error: ", err)

		case <-time.After(10 * time.Second):
			fmt.Println("Sorry time is up!")

			err := ResetRepository(repositoryPath)
			if err != nil {
				fmt.Println("You got lucky, there was an error while ressetting: ", err)
			}

		case cancelled = <-cancel:
			fmt.Println("\nBye!")
		}
	}
}
