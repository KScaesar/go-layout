package dataflow

import (
	"context"
	"sync"
)

// Return 表示 Sender 發送 Message 後, Receiver 處理後的結果.
// Receiver 透過 "Async" 發送 Return.
// Sender   透過 "Await" 接收 Return.
//
// https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html
// https://docs.nats.io/nats-concepts/core-nats/reqreply#the-pattern
type Return struct {
	Result any
	Err    error
}

func NewReply(qty int) Reply {
	return Reply{
		qty: qty,
		mq:  make(chan *Return, qty),
	}
}

type Reply struct {
	qty int
	mq  chan *Return
}

func (r Reply) Async(Result any, Err error) {
	r.AsyncWithCtx(context.Background(), Result, Err)
}

func (r Reply) Await() (Result any, Err error) {
	return r.AwaitWithCtx(context.Background())
}

func (r Reply) MultiAwait() (Results []any, Err error) {
	return r.MultiAwaitWithCtx(context.Background())
}

func (r Reply) AsyncWithCtx(ctx context.Context, Result any, Err error) (err error) {
	_return := &Return{
		Result: Result,
		Err:    Err,
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case r.mq <- _return:
		err = nil
	}
	return err
}

func (r Reply) AwaitWithCtx(ctx context.Context) (Result any, Err error) {
	select {
	case <-ctx.Done():
		Result = nil
		Err = ctx.Err()
	case _return := <-r.mq:
		Result = _return.Result
		Err = _return.Err
	}
	return
}

func (r Reply) MultiAwaitWithCtx(ctx context.Context) (Results []any, Err error) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	Results = make([]any, 0, r.qty)
	done := make(chan error, r.qty)

	for i := 0; i < r.qty; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := r.AwaitWithCtx(ctx)
			if err != nil {
				done <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()
			Results = append(Results, result)
		}()
	}

	go func() {
		wg.Wait()
		done <- nil
	}()

	if err := <-done; err != nil {
		return []any{}, err
	}
	return Results, nil
}

//

func MultiAwait(multiReply []Reply) (Results []any, Err error) {
	return MultiAwaitWithCtx(context.Background(), multiReply)
}

func MultiAwaitWithCtx(ctx context.Context, multiReply []Reply) (Results []any, Err error) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	qty := len(multiReply)
	Results = make([]any, qty)
	done := make(chan error, qty)

	for i, reply := range multiReply {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			result, err := reply.AwaitWithCtx(ctx)
			if err != nil {
				done <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()
			Results[idx] = result
		}(i)
	}

	go func() {
		wg.Wait()
		done <- nil
	}()

	if err := <-done; err != nil {
		return []any{}, err
	}
	return Results, nil
}
