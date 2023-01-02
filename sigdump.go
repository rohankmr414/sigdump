package sigdump

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"golang.org/x/sys/unix"
)

func Start() {
	signalStr := os.Getenv("SIGDUMP_SIGNAL")
	if signalStr == "" {
		signalStr = "SIGCONT"
	}

	s := unix.SignalNum(signalStr)
	if s == 0 {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, s)

	go func() {
		for {
			s := <-c
			dumpStack(s)
		}
	}()
}

func dumpStack(s os.Signal) {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)

	path := os.Getenv("SIGDUMP_PATH")
	if path == "" {
		path = fmt.Sprintf("/tmp/sigdump-%d.log", os.Getpid())
	}

	var w *os.File
	if path == "-" {
		w = os.Stdout
	} else if path == "+" {
		w = os.Stderr
	} else {
		var err error
		w, err = os.Create(path)
		if err != nil {
			return
		}
		defer w.Close()
	}

	fmt.Fprintf(w, "Received signal %s\n%s", s, buf[:n])
}
