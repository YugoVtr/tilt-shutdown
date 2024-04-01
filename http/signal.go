package http

import (
	"os"
	"os/signal"
	"syscall"
)

func newInterruptSignalChannel() <-chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	return stop
}
