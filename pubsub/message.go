package pubsub

// base, all message has one base which show the basic of a message
type base struct {
	topic string
	err   error
}

func (b *base) Topic() string {
	return b.topic
}

func (b *base) Error() error {
	return b.err
}

type message struct {
	base
	msg interface{}
}

func FromMessage(topic string, msg interface{}) *message {
	return &message{base: base{topic: topic}, msg: msg}
}

func (m *message) Message() interface{} {
	return m.msg
}

func (m *message) Unmarshal(v interface{}) error {
	panic("implement me !")
}
