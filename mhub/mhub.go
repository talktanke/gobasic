package mhub

import (
	"context"
)

// MHub a hub use to publish a series of message, in different of medium like memory or tcp....
type MHub interface {
	// Publish a message so that the hub will push these message one by one to subscriber.
	Publish(ctx context.Context, topic string, message interface{})

	// Subscribe input a topic and get a handle(Subscription) to receive message
	Subscribe(ctx context.Context, topic string) Subscription
}

// Subscription is return by Subscribe, like a handle to receive the message or close it.
type Subscription interface {
	// Close should implement to wait all message send over
	Close()
	// Chan returns a <-chan like a handle to read message
	Chan() <-chan Message
}

// Message interface abstract what the msg do
type Message interface {
	Topic() string
	Message() interface{}
	Error() error
	Unmarshal(v interface{}) error
}
