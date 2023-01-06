package sigdump

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

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

	now := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	ppid := os.Getppid()
	_, err := fmt.Fprintf(w, "time=%s host=%s pid=%d ppid=%d signal=%s\n%s", now, hostname, pid, ppid, s, buf[:n])
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "\nmemory profile:\n")
	if err != nil {
		return
	}
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	_, err = fmt.Fprintf(w, "Alloc = %v\n", memStats.Alloc)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "TotalAlloc = %v\n", memStats.TotalAlloc)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Sys = %v\n", memStats.Sys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Lookups = %v\n", memStats.Lookups)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Mallocs = %v\n", memStats.Mallocs)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Frees = %v\n", memStats.Frees)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapAlloc = %v\n", memStats.HeapAlloc)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapSys = %v\n", memStats.HeapSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapIdle = %v\n", memStats.HeapIdle)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapInuse = %v\n", memStats.HeapInuse)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapReleased = %v\n", memStats.HeapReleased)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "HeapObjects = %v\n", memStats.HeapObjects)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "StackInuse = %v\n", memStats.StackInuse)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "StackSys = %v\n", memStats.StackSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "MSpanInuse = %v\n", memStats.MSpanInuse)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "MSpanSys = %v\n", memStats.MSpanSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "MCacheInuse = %v\n", memStats.MCacheInuse)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "MCacheSys = %v\n", memStats.MCacheSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "BuckHashSys = %v\n", memStats.BuckHashSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "GCSys = %v\n", memStats.GCSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "OtherSys = %v\n", memStats.OtherSys)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "NextGC = %v\n", memStats.NextGC)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "LastGC = %v\n", memStats.LastGC)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "PauseTotalNs = %v\n", memStats.PauseTotalNs)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "NumGC = %v\n", memStats.NumGC)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "GCCPUFraction = %v\n", memStats.GCCPUFraction)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "DebugGC = %v\n", memStats.DebugGC)
	if err != nil {
		return
	}
}
