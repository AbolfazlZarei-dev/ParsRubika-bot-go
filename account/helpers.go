package account

import (
	"bufio"
	"os"
	"sync"
	"time"
)

// stdin returns the singleton bufio.Reader wrapping os.Stdin. The
// reader is initialised lazily on first use so importing the account
// package never touches stdin (important for tests and headless use).
var (
	stdinOnce   sync.Once
	stdinReader *bufio.Reader
)

func stdin() *bufio.Reader {
	stdinOnce.Do(func() {
		stdinReader = bufio.NewReader(os.Stdin)
	})
	return stdinReader
}

// durationOfSeconds returns a time.Duration of `n` seconds.
func durationOfSeconds(n int) time.Duration { return time.Duration(n) * time.Second }

// timerAfterSeconds returns a channel that fires after `n` seconds.
func timerAfterSeconds(n int) <-chan time.Time { return time.After(time.Duration(n) * time.Second) }
