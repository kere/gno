package util

import (
	"os"
	"os/signal"
	"syscall"
)

// ListenSignal os signal
func ListenSignal(f func(sign os.Signal)) {
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		f(<-quitChan)
	}()
}
