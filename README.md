# sigdump

`sigdump` is a Go package that captures and dumps the stack trace to a file when a specified signal is received. It is inspired by the [fluent/sigdump](https://github.com/fluent/sigdump) gem for Ruby.

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

By default, `sigdump` will listen for the `SIGCONT` signal and print the stack trace to a file in the `/tmp` directory when the signal is received. To listen for a different signal, set the `SIGDUMP_SIGNAL` environment variable to the signal name you want to listen for. For example, to listen for the SIGINT signal (Ctrl+C), you can set the environment variable like this:

```bash
$ SIGDUMP_SIGNAL=SIGINT myprogram
```

To specify a different path to write the stack trace to, set the `SIGDUMP_PATH` environment variable to the desired path. The special values "-" and "+" can be used to write the stack trace to `stdout` and `stderr`, respectively.

For example, to write the stack trace to a file in the `/var/log` directory, you can set the `SIGDUMP_PATH` environment variable like this:

```bash
$ SIGDUMP_PATH=/var/log/sigdump.log myprogram
```
