package sigdump

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Start sets up a signal handler for the signal specified by the SIGDUMP_SIGNAL
// environment variable. If the environment variable is not set, the default signal is SIGCONT.
// If the signal is received, a runtime stack trace and memory profile of the process
// will be written to the file specified by the SIGDUMP_PATH environment variable.
// If the environment variable is not set, the default file is /tmp/sigdump-<pid>.log.
// If the value of SIGDUMP_PATH is - the stack trace is written to the stdout.
// If the value of SIGDUMP_PATH is + the stack trace is written to the stderr.
// The stack trace includes the current time, hostname, pid, ppid and signal.
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

	buf = []byte(strings.ReplaceAll("  "+string(buf[:n]), "\n", "\n\t"))
	buf = []byte(strings.ReplaceAll(string(buf), "\tgoroutine ", "  goroutine "))
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
		defer func(w *os.File) {
			err = w.Close()
			if err != nil {
				return
			}
		}(w)
	}

	now := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	ppid := os.Getppid()
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("Sigdump time=%s host=%s pid=%d ppid=%d signal=%s\n%s\n", now, hostname, pid, ppid, s, buf[:n]))

	sb.WriteString("\n  Mem Stat:\n")
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	sb.WriteString(fmt.Sprintf("\tAlloc = %v\n", memStats.Alloc))
	sb.WriteString(fmt.Sprintf("\tTotalAlloc = %v\n", memStats.TotalAlloc))
	sb.WriteString(fmt.Sprintf("\tSys = %v\n", memStats.Sys))
	sb.WriteString(fmt.Sprintf("\tLookups = %v\n", memStats.Lookups))
	sb.WriteString(fmt.Sprintf("\tMallocs = %v\n", memStats.Mallocs))
	sb.WriteString(fmt.Sprintf("\tFrees = %v\n", memStats.Frees))
	sb.WriteString(fmt.Sprintf("\tHeapAlloc = %v\n", memStats.HeapAlloc))
	sb.WriteString(fmt.Sprintf("\tHeapSys = %v\n", memStats.HeapSys))
	sb.WriteString(fmt.Sprintf("\tHeapIdle = %v\n", memStats.HeapIdle))
	sb.WriteString(fmt.Sprintf("\tHeapInuse = %v\n", memStats.HeapInuse))
	sb.WriteString(fmt.Sprintf("\tHeapReleased = %v\n", memStats.HeapReleased))
	sb.WriteString(fmt.Sprintf("\tHeapObjects = %v\n", memStats.HeapObjects))
	sb.WriteString(fmt.Sprintf("\tStackInuse = %v\n", memStats.StackInuse))
	sb.WriteString(fmt.Sprintf("\tStackSys = %v\n", memStats.StackSys))
	sb.WriteString(fmt.Sprintf("\tMSpanInuse = %v\n", memStats.MSpanInuse))
	sb.WriteString(fmt.Sprintf("\tMSpanSys = %v\n", memStats.MSpanSys))
	sb.WriteString(fmt.Sprintf("\tMCacheInuse = %v\n", memStats.MCacheInuse))
	sb.WriteString(fmt.Sprintf("\tMCacheSys = %v\n", memStats.MCacheSys))
	sb.WriteString(fmt.Sprintf("\tBuckHashSys = %v\n", memStats.BuckHashSys))
	sb.WriteString(fmt.Sprintf("\tGCSys = %v\n", memStats.GCSys))
	sb.WriteString(fmt.Sprintf("\tOtherSys = %v\n", memStats.OtherSys))
	sb.WriteString(fmt.Sprintf("\tNextGC = %v\n", memStats.NextGC))
	sb.WriteString(fmt.Sprintf("\tLastGC = %v\n", memStats.LastGC))
	sb.WriteString(fmt.Sprintf("\tPauseTotalNs = %v\n", memStats.PauseTotalNs))
	sb.WriteString(fmt.Sprintf("\tNumGC = %v\n", memStats.NumGC))
	sb.WriteString(fmt.Sprintf("\tGCCPUFraction = %v\n", memStats.GCCPUFraction))
	sb.WriteString(fmt.Sprintf("\tDebugGC = %v\n", memStats.DebugGC))

	_, err := w.Write([]byte(sb.String()))
	if err != nil {
		return
	}
}
