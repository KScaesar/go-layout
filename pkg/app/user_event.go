package app

import (
	"github.com/KScaesar/go-layout/pkg/utility/dataflow"
)

func NewRegisteredUserEvent(user *User) *dataflow.Message {
	return dataflow.NewBodyEgress("user.registered", &RegisteredUserEvent{})
}

type RegisteredUserEvent struct {
}
