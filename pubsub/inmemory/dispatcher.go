package inmemory

import (
	"context"
	"sync"

	"github.com/talktanke/gobasic/log"
	"github.com/talktanke/gobasic/pubsub"


)

// Dispatcher is in charge of a topic
type Dispatcher struct {
	chans []chan<- pubsub.Message
	sync.RWMutex
}

// NewDispatcher create new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{chans: make([]chan<- pubsub.Message, 0, 0)}
}

// Add one subscriber to the topic
func (d *Dispatcher) Add(ch chan<- pubsub.Message) {
	d.Lock()
	defer d.Unlock()
	d.chans = append(d.chans, ch)
}

// Remove removes a chan from chans
func (d *Dispatcher) Remove(ch chan pubsub.Message) {
	d.Lock()
	defer d.Unlock()
	nch := make([]chan<- pubsub.Message, len(ch))
	for i := 0; i < len(d.chans); i++ {
		if ch != d.chans[i] {
			nch = append(nch, d.chans[i])
		}
	}
	d.chans = nch
}

func (d *Dispatcher) send(ch chan<- pubsub.Message, msg pubsub.Message) {
	select {
	case ch <- msg:
		return
	default:
		log.Debugf("message send fail, channel blocked.")
	}
}

// putOneByOne like it's name, put message one to subscriber one by one
func (d *Dispatcher) putOneByOne(ctx context.Context, message pubsub.Message) {
	d.RLock()
	chans := d.chans
	d.RUnlock()
	for _, ch := range chans {
		select {
		case <-ctx.Done():
			return
		case ch <- message:
			// log.Debugf("message %v has been send to %s", message.Message(), name)
		}
	}
}
