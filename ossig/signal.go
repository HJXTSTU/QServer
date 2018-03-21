package ossig

import (
	"os"
	"sync"
	"os/signal"
	"syscall"
)


var signal_exit chan os.Signal

var one sync.Once

func GetExitChan()<-chan os.Signal{
	one.Do(func() {
		signal_exit = make(chan os.Signal)
		signal.Notify(signal_exit,syscall.SIGINT,syscall.SIGKILL,syscall.SIGTERM)
	})
	return signal_exit
}