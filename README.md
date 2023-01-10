# sigdump

`sigdump` is a Go package that captures and dumps the stack trace and memory profile to a file when a specified signal
is received. It is inspired by the [fluent/sigdump](https://github.com/fluent/sigdump) gem for Ruby.

## Installation

```bash
$ go get github.com/rohankmr414/sigdump
```

## Usage

To use `sigdump`, import the package and call the Start function at the beginning of your program:

```bash
import "github.com/rohankmr414/sigdump"

func main() {
	sigdump.Start()
	// Rest of your program goes here...
}
```

By default, `sigdump` will listen for the `SIGCONT` signal and print the stack trace to a file in the `/tmp` directory
when the signal is received. To listen for a different signal, set the `SIGDUMP_SIGNAL` environment variable to the
signal name you want to listen for. For example, to listen for the SIGINT signal (Ctrl+C), you can set the environment
variable like this:

```bash
$ SIGDUMP_SIGNAL=SIGINT myprogram
```

To specify a different path to write the stack trace to, set the `SIGDUMP_PATH` environment variable to the desired
path. The special values "-" and "+" can be used to write the stack trace to `stdout` and `stderr`, respectively.

For example, to write the stack trace to a file in the `/var/log` directory, you can set the `SIGDUMP_PATH` environment
variable like this:

```bash
$ SIGDUMP_PATH=/var/log/sigdump.log myprogram
```

## Sample Output

```
$ cat /tmp/sigdump-8829.log
Sigdump time=2023-01-10T18:01:39+05:30 host=MacBook.local pid=8829 ppid=5521 signal=SIGCONT
  goroutine 34 [running]:
	github.com/rohankmr414/sigdump.dumpStack({0x1002ed918?, 0x10035bcd8})
		/Users/rohan/Documents/code/personal/go-test/vendor/github.com/rohankmr414/sigdump/sigdump.go:46 +0x68
	github.com/rohankmr414/sigdump.Start.func1()
		/Users/rohan/Documents/code/personal/go-test/vendor/github.com/rohankmr414/sigdump/sigdump.go:39 +0x38
	created by github.com/rohankmr414/sigdump.Start
		/Users/rohan/Documents/code/personal/go-test/vendor/github.com/rohankmr414/sigdump/sigdump.go:36 +0xd4

  goroutine 1 [sleep]:
	time.Sleep(0x1bf08eb000)
		/usr/local/go/src/runtime/time.go:195 +0x118
	main.main()
		/Users/rohan/Documents/code/personal/go-test/main.go:12 +0x2c

  goroutine 33 [syscall]:
	os/signal.signal_recv()
		/usr/local/go/src/runtime/sigqueue.go:149 +0x2c
	os/signal.loop()
		/usr/local/go/src/os/signal/signal_unix.go:23 +0x1c
	created by os/signal.Notify.func1.1
		/usr/local/go/src/os/sig

  Mem Stat:
	Alloc = 203440
	TotalAlloc = 203440
	Sys = 8834064
	Lookups = 0
	Mallocs = 258
	Frees = 9
	HeapAlloc = 203440
	HeapSys = 3801088
	HeapIdle = 3096576
	HeapInuse = 704512
	HeapReleased = 3063808
	HeapObjects = 249
	StackInuse = 393216
	StackSys = 393216
	MSpanInuse = 37536
	MSpanSys = 48960
	MCacheInuse = 9600
	MCacheSys = 15600
	BuckHashSys = 3523
	GCSys = 3776672
	OtherSys = 795005
	NextGC = 4194304
	LastGC = 0
	PauseTotalNs = 0
	NumGC = 0
	GCCPUFraction = 0
	DebugGC = false
```
