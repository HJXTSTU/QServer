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

var goroutines_wait sync.WaitGroup

var wait_one sync.Once

func GetGlobalWaitGroup()sync.WaitGroup{
	one.Do(func() {
		goroutines_wait= sync.WaitGroup{}
	})
	return goroutines_wait
}