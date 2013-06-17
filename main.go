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
	builder := NewTimelineBuilder(1 * time.Second)
	builder.AddEvent(time.Duration(2*time.Second), func(tickEvent *TickEvent) {
		fmt.Println("2 seconds tick event!")
	})
	builder.AddEvent(time.Duration(3*time.Second), func(tickEvent *TickEvent) {
		fmt.Println("3 seconds tick event!")
	})

	timeline := builder.Build()
	timeline.Start()

	for !cancelled {
		select {
		case newHeadCommit := <-watcher.HeadChanged:
			fmt.Println("new head: ", newHeadCommit.Id())
			timeline.Reset()

		case err := <-watcher.Error:
			fmt.Println("error: ", err)

		case <-timeline.Tick:
			fmt.Println("tick...")

		case <-timeline.Finished:
			fmt.Println("Sorry time is up!")

			//err := ResetRepository(repositoryPath)
			//if err != nil {
			//	fmt.Println("You got lucky, there was an error while resetting: ", err)
			//}
			timeline.Start()

		case cancelled = <-cancel:
			fmt.Println("\nBye!")
		}
	}
}
