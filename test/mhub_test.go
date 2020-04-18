package test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dashjay/gobasic/mhub/inmemory"
)

func TestInMemory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inm := inmemory.New()
	sub := inm.Subscribe(ctx, "love", "dj")
	ch := sub.Chan()
	inm.Publish(ctx, "love", 1)
	msg := <-ch
	require.Equal(t, 1, msg.Message())
	inm.Publish(ctx, "love", 2)
	msg = <-ch
	require.Equal(t, 2, msg.Message())
	cancel()
	msg, more := <-ch
	require.Equal(t, false, more)
}

func TestInMemoryMultiSubscriber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inm := inmemory.New()
	dj1 := inm.Subscribe(ctx, "love", "dj1")
	dj2 := inm.Subscribe(ctx, "love", "dj2")
	ch1 := dj1.Chan()
	ch2 := dj2.Chan()
	inm.Publish(ctx, "love", 666)
	m1 := <-ch1
	t.Logf("dj1 receive message %v", m1.Message())
	m2 := <-ch2
	t.Logf("dj2 receive message %v", m2.Message())
	cancel()
}

func BenchmarkInMemory(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	inm := inmemory.New()
	dj1 := inm.Subscribe(ctx, "love", "dj1")
	dj2 := inm.Subscribe(ctx, "love", "dj2")

	done := make(chan struct{})
	go func() {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inm.Publish(ctx, "love", 666)
		}
		done <- struct{}{}
		b.StopTimer()
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range dj1.Chan() {
		}
	}()
	go func() {
		defer wg.Done()
		for range dj2.Chan() {
		}
	}()

	<-done
	cancel()
	wg.Wait()
}
