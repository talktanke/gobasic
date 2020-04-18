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

// Subscribe subscribe one topic return a instance of subscription for getting message.
func (i *inMemory) Subscribe(ctx context.Context, option *mhub.SubscriptionOptions) mhub.Subscription {
	if option.BufferSize < 0 {
		option.BufferSize = 0
	}
	// TODO: buffer should be a parameter.....
	ch := make(chan mhub.Message, option.BufferSize)
	ready := make(chan struct{})
	ctx, sub := newSubscription(ctx, ch)
	go i.watch(ctx, option.Topics, sub, ready)
	<-ready
	return sub
}

func (i *inMemory) watch(ctx context.Context, topics []string, sub *subscription, ready chan<- struct{}) {
	defer func() {
		sub.cancel()
		log.Debugf("stop subscribe topic %v", topics)
		for _, t := range topics {
			d := i.getDispatcher(t)
			d.Remove(sub.ch)
		}
		close(sub.ch)
	}()
	for _, t := range topics {
		d := i.getDispatcher(t)
		d.Add(sub.ch)
	}
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
