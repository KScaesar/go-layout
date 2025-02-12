package dataflow

import (
	"context"
	"sync"
)

type response struct {
	Result any
	Err    error
}

func NewReply(qty int) Reply {
	return Reply{
		qty: qty,
		mq:  make(chan *response, qty),
	}
}

// Reply is used to push or pull a response.
// The response represents the result obtained after the Consumer processes the Message send by the Producer.
//
// - Consumer use Reply.Push response.
//
// - Producer use Reply.Pull response.
//
// https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html
// https://docs.nats.io/nats-concepts/core-nats/reqreply#the-pattern
type Reply struct {
	qty int
	mq  chan *response
}

func (r Reply) Push(Result any, Err error) {
	r.PushWithCtx(context.Background(), Result, Err)
}

func (r Reply) PushWithCtx(ctx context.Context, Result any, Err error) (err error) {
	_return := &response{
		Result: Result,
		Err:    Err,
	}

	select {
	case <-ctx.Done():
		err = context.Cause(ctx)
	case r.mq <- _return:
		err = nil
	}
	return err
}

// Pull 代表的意思是, Producer 透過 一個 Reply 取得 一個 Message 的 response
func (r Reply) Pull() (Result any, Err error) {
	return r.PullWithCtx(context.Background())
}

func (r Reply) PullWithCtx(ctx context.Context) (Result any, Err error) {
	select {
	case <-ctx.Done():
		Result = nil
		Err = context.Cause(ctx)
	case _return := <-r.mq:
		Result = _return.Result
		Err = _return.Err
	}
	return
}

// Pulls 代表的意思是, Producer 透過 一個 Reply 取得 多個 Message 的 response
func (r Reply) Pulls() (Results []any, Err error) {
	return r.PullsWithCtx(context.Background())
}

func (r Reply) PullsWithCtx(ctx context.Context) (Results []any, Err error) {
	Results = make([]any, 0, r.qty)

	for i := 0; i < r.qty; i++ {
		result, err := r.PullWithCtx(ctx)
		if err != nil {
			return []any{}, err
		}
		Results = append(Results, result)
	}

	return Results, nil
}

//

// Gather 代表的意思是, Producer 透過 多個 Reply 取得 多個 Message 的 response, 其中 Message 和 Reply 是 1 對 1 關聯
func Gather(multiReply []Reply) (Results []any, Err error) {
	return GatherWithCtx(context.Background(), multiReply)
}

func GatherWithCtx(ctx context.Context, multiReply []Reply) (Results []any, Err error) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	qty := len(multiReply)
	Results = make([]any, qty)
	ch := make(chan error)

	ctx, cancel := context.WithCancelCause(ctx)

	for idx, reply := range multiReply {
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := reply.PullWithCtx(ctx)
			if err != nil {
				select {
				case ch <- err:
				default:
				}
				return
			}

			mu.Lock()
			Results[idx] = result
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		select {
		case ch <- nil:
		default:
		}
	}()

	if err := <-ch; err != nil {
		cancel(err)
		return []any{}, err
	}
	return Results, nil
}
