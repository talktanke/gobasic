package inmemory

import (
	"context"
	"sync"

	"github.com/dashjay/gobasic/log"
	"github.com/dashjay/gobasic/mhub"
)

// subscription will be hold by subscriber to get the published message
type subscription struct {
	ch     chan mhub.Message
	cancel context.CancelFunc
}

// newSubscription new a  subscription
func newSubscription(ctx context.Context, ch chan mhub.Message) (context.Context, *subscription) {
	subCtx, cancel := context.WithCancel(ctx)
	return subCtx, &subscription{ch: ch, cancel: cancel}
}

func (s *subscription) Chan() <-chan mhub.Message {
	return s.ch
}

func (s *subscription) Close() {
	s.cancel()
}

// inMemory, a hub in memory manages a lot topic
type inMemory struct {
	topics map[string]*Dispatcher
	sync.Mutex
}

func New() *inMemory {
	return &inMemory{
		topics: make(map[string]*Dispatcher),
	}
}

// Publish publishes a message to the topic
func (i *inMemory) Publish(ctx context.Context, topic string, msg interface{}) {
	d := i.getDispatcher(topic)
	sub := mhub.FromMessage(topic, msg)
	d.putOneByOne(ctx, sub)
}

func (i *inMemory) Subscribe(ctx context.Context, topic, name string) mhub.Subscription {
	// TODO: buffer should be a parameter.....
	ch := make(chan mhub.Message, 64)
	ready := make(chan struct{})
	ctx, sub := newSubscription(ctx, ch)
	go i.watch(ctx, topic, name, sub, ready)
	<-ready
	return sub
}

func (i *inMemory) watch(ctx context.Context, topic, name string, sub *subscription, ready chan<- struct{}) {
	defer func() {
		sub.cancel()
		log.Debugf("%s stop subscribe topic %s", name, topic)
		d := i.getDispatcher(topic)
		d.Remove(name)
		close(sub.ch)
	}()
	d := i.getDispatcher(topic)
	d.Add(sub.ch, name)
	ready <- struct{}{}
	<-ctx.Done()
}

// getDispatcher get dispatcher by topic
func (i *inMemory) getDispatcher(topic string) *Dispatcher {
	i.Lock()
	defer i.Unlock()
	d, ok := i.topics[topic]
	if !ok {
		d = NewDispatcher()
		i.topics[topic] = d
	}
	return d
}
