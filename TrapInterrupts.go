package main

import (
	"os"
	"os/signal"
	"syscall"
)

func TrapInterrupts() <-chan (bool) {
	interrupt_chan := make(chan os.Signal, 2)
	cancel := make(chan bool)

	signal.Notify(interrupt_chan, os.Interrupt)
	signal.Notify(interrupt_chan, syscall.SIGTERM)

	go func() {
		<-interrupt_chan
		cancel <- true
	}()

	return cancel
}
