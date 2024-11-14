package dataflow

import (
	"context"
	"sync"

	"github.com/gookit/goutil/maputil"

	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewBodyEgress(subject string, body any) *Message {
	egress := GetMessage()

	egress.Subject = subject
	egress.Body = body
	return egress
}

func NewBytesEgress(subject string, bMessage []byte) *Message {
	egress := GetMessage()

	egress.Subject = subject
	egress.Bytes = bMessage
	return egress
}

func GetMessage() *Message {
	return messagePool.Get()
}

func PutMessage(message *Message) {
	message.reset()
	messagePool.Put(message)
}

var messagePool = utility.NewPool(newMessage)

func newMessage() *Message {
	return &Message{
		RouteParam: map[string]any{},
		Metadata:   map[string]any{},
		Ctx:        context.Background(),
	}
}

// Message represents a high-level abstraction data structure containing metadata (e.g. header) + body
type Message struct {
	Subject string

	Bytes []byte // ingress byte payload or egress byte payload
	Body  any    // egress golang object

	identifier string

	Mutex sync.Mutex // for egress broadcast data process

	// RouteParam are used to capture values from subject.
	// These parameters represent resources or identifiers.
	//
	// Example:
	//
	//	define mux subject = "/users/{id}"
	//	send or recv subject = "/users/1017"
	//
	//	get route param:
	//		key : value => id : 1017
	RouteParam maputil.Data

	Metadata maputil.Data

	// raw message from 3rd pkg
	//
	// Example:
	//
	// fiber.Ctx or amqp.Delivery or kafka.Message
	RawInfra any

	reply    Reply
	pingpong chan struct{} // websocket or tcp socket for check connection health

	Ctx context.Context
}

func (msg *Message) MsgId() string {
	if msg.identifier == "" {
		msg.identifier = utility.NewUlid()
	}
	return msg.identifier
}

func (msg *Message) SetMsgId(msgId string) {
	msg.identifier = msgId
}

func (msg *Message) reset() {
	msg.Subject = ""
	msg.Bytes = nil
	msg.Body = nil
	msg.identifier = ""

	for key := range msg.RouteParam {
		delete(msg.RouteParam, key)
	}
	for key := range msg.Metadata {
		delete(msg.Metadata, key)
	}

	msg.RawInfra = nil
	msg.reply.mq = nil
	msg.pingpong = nil

	msg.Ctx = context.Background()
}

func (msg *Message) Copy() *Message {
	message := GetMessage()

	message.Subject = msg.Subject
	message.Bytes = msg.Bytes
	message.Body = msg.Body
	message.identifier = msg.identifier

	for key, v := range msg.RouteParam {
		message.RouteParam.Set(key, v)
	}
	for key, v := range msg.Metadata {
		message.Metadata.Set(key, v)
	}

	message.RawInfra = msg.RawInfra
	message.reply = msg.reply
	message.pingpong = msg.pingpong

	message.Ctx = msg.Ctx
	return message
}

func (msg *Message) Reply() Reply {
	if msg.reply.mq == nil {
		msg.reply = NewReply(1)
	}
	return msg.reply
}

func (msg *Message) SetReply(r Reply) {
	msg.reply = r
}

func (msg *Message) AckPingPong() {
	if msg.pingpong == nil {
		panic("pingpong channel is nil")
	}
	msg.pingpong <- struct{}{}
}

func (msg *Message) SetPingPong(pingpong chan struct{}) {
	msg.pingpong = pingpong
}
