package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// cancel channel, writing to this
	// will signal the monitoring operation
	// to cancel.
	cancel := make(chan bool, 1)

	// wire (SIGINT and SIGTERM) to
	// write to the cancel channel, this allows
	// the program to react on [ctrl]+[c] key
	// presses from the console.
	trapSignal(cancel)

	run(cancel)
}

func run(cancel chan (bool)) {
	monitor(cancel, 1*time.Second, func() {
		fmt.Println(time.Now)
	})
}

func monitor(cancel chan (bool), duration time.Duration, work func()) {
	cancelled := false
	for !cancelled {
		timeout := time.After(duration)
		select {
		case <-cancel:
			fmt.Printf("quiting from work...")
			cancelled = true
		case <-timeout:
			work()
		}
	}

	fmt.Println("exiting prommise...")
}

func work() {
	fmt.Printf("%v", time.Now)
}

func trapSignal(cancel chan (bool)) {
	interrupt_chan := make(chan os.Signal, 2)

	signal.Notify(interrupt_chan, os.Interrupt)
	signal.Notify(interrupt_chan, syscall.SIGTERM)

	fmt.Println("started... waiting for SIGTERM")

	go func() {
		sig := <-interrupt_chan
		fmt.Printf("\nsignal %s received, cancelling...", sig)
		cancel <- true
	}()
}
