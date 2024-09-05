package dataflow

type MessageBus interface {

	// Send: Messages are processed and sent.
	//
	// Before transmission, messages are subjected to a data processing flow, such as
	// encoding or serialization, and then delivered to an external medium.
	//
	// Example:
	//	egress1 := &Message{Subject: "user.loggedin", Body: map[string]any{"UserID": "123"}
	//	egress2 := &Message{Subject: "user.purchased", Body: map[string]any{"UserID": "456"}
	//	MessageBus.Send(egress1, egress2)
	Send(messages ...*Message) error

	// RawSend: Messages are directly sent without any processing.
	//
	// The messages are transmitted in their raw form, without encoding or
	// serialization, and are delivered directly to the external medium.
	//
	// Example:
	//	egress3 := &Message{Subject: "event3", Bytes: []byte("rawdata3")}
	//	egress4 := &Message{Subject: "event4", Bytes: []byte("rawdata4")}
	//	MessageBus.RawSend(egress3, egress4)
	RawSend(messages ...*Message) error
}
