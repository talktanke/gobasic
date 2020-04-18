package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dashjay/gobasic/time/backoff"
)

func TestBackOff(t *testing.T) {
	b := &backoff.Backoff{
		// These are the defaults
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	fmt.Printf("%s\n", b.Duration())
	fmt.Printf("%s\n", b.Duration())
	fmt.Printf("%s\n", b.Duration())

	fmt.Printf("Reset!\n")
	b.Reset()

	fmt.Printf("%s\n", b.Duration())
}
