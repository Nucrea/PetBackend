package request_context

import (
	"context"
	"time"
)

type RequestContext interface {
	context.Context

	RequestId() string
	UserId() string
}

type Opts struct {
	RequestId string
	UserId    string
}

func New(opts Opts) RequestContext {
	return reqCtx{
		requestId: opts.RequestId,
		userId:    opts.UserId,
	}
}

type reqCtx struct {
	ctx       context.Context
	requestId string
	userId    string
}

func (r reqCtx) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

func (r reqCtx) Done() <-chan struct{} {
	return r.ctx.Done()
}

func (r reqCtx) Err() error {
	return r.ctx.Err()
}

func (r reqCtx) Value(key any) any {
	return r.ctx.Value(key)
}

func (r reqCtx) String() string {
	return r.ctx.Err().Error()
}

func (r reqCtx) RequestId() string {
	return r.requestId
}

func (r reqCtx) UserId() string {
	return r.userId
}
