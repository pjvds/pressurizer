package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cancel := trapSignal()

	run(cancel)
}

func run(cancel <-chan (bool)) {
	stop := false
	for !stop {
		fmt.Println(time.Now())

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
