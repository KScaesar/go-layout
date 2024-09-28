package dataflow

import (
	"context"
)

//go:generate mockgen -typed -package=dataflow -destination=pubsub_mock.go github.com/KScaesar/go-layout/pkg/utility/dataflow Producer
type Producer interface {
	Send(messages ...*Message) error
	SendWithCtx(ctx context.Context, messages ...*Message) error
}

type Consumer interface {
	Listen() (err error)
	Stop() error
}
