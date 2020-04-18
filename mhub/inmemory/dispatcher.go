package inmemory

import (
	"context"
	"sync"

	"github.com/dashjay/gobasic/log"
	"github.com/dashjay/gobasic/mhub"
)

// Dispatcher is in charge of a topic
type Dispatcher struct {
	chans []chan<- mhub.Message
	sync.RWMutex
}

// NewDispatcher create new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{chans: make([]chan<- mhub.Message, 0, 0)}
}

// Add one subscriber to the topic
func (d *Dispatcher) Add(ch chan<- mhub.Message) {
	d.Lock()
	defer d.Unlock()
	d.chans = append(d.chans, ch)
}

// Remove removes a chan from chans
func (d *Dispatcher) Remove(ch chan mhub.Message) {
	d.Lock()
	defer d.Unlock()
	nch := make([]chan<- mhub.Message, len(ch))
	for i := 0; i < len(d.chans); i++ {
		if ch != d.chans[i] {
			nch = append(nch, d.chans[i])
		}
	}
	d.chans = nch
}

func (d *Dispatcher) send(ch chan<- mhub.Message, msg mhub.Message) {
	select {
	case ch <- msg:
		return
	default:
		log.Debugf("message send fail, channel blocked.")
	}
}

// putOneByOne like it's name, put message one to subscriber one by one
func (d *Dispatcher) putOneByOne(ctx context.Context, message mhub.Message) {
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
