package pubsub

// subscriptionOptions can help subscriber do subscribe on more than one topic
// and use buffer size as a parameter
type SubscriptionOptions struct {
	BufferSize int
	Topics     []string
}

func NewSubscriptionOptions(bufSize int, topics ...string) *SubscriptionOptions {
	return &SubscriptionOptions{BufferSize: bufSize, Topics: topics}
}
