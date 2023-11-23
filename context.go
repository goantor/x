package x

import (
	"context"
	"time"
)

type IContextData interface {
	GiveService(service string)
	TakeService() string
	GiveModule(module string)
	TakeModule() string
	GiveAction(action string)
	TakeAction() string
	TakeTraceId() string
	GiveTraceId(id string)
	TakeRequestId() string
	GiveRequestId(id string)
	GiveUser(user interface{})
	GiveParams(params interface{})
	GiveIP(ip string)
	Set(key string, value interface{})
	Get(key string, def interface{}) interface{}
}

func NewContextData() IContextData {
	return &ContextData{
		Data: make(H),
	}
}

type ContextData struct {
	Service   string      `json:"service,omitempty"`
	Module    string      `json:"module,omitempty"`
	Action    string      `json:"action,omitempty"`
	TraceId   string      `json:"trace_id,omitempty"`
	RequestId string      `json:"request_id,omitempty"`
	User      interface{} `json:"user,omitempty"`
	Params    interface{} `json:"params,omitempty"`
	IP        string      `json:"ip,omitempty"`
	Data      H           `json:"data,omitempty"`
}

func (c *ContextData) GiveService(service string) {
	c.Service = service
}

func (c *ContextData) TakeService() string {
	return c.Service
}

func (c *ContextData) GiveModule(module string) {
	c.Module = module
}

func (c *ContextData) TakeModule() string {
	return c.Module
}

func (c *ContextData) GiveAction(action string) {
	c.Action = action
}

func (c *ContextData) TakeAction() string {
	return c.Action
}

func (c *ContextData) TakeTraceId() string {
	return c.TraceId
}

func (c *ContextData) GiveRequestId(id string) {
	c.RequestId = id
}

func (c *ContextData) TakeRequestId() string {
	return c.RequestId
}

func (c *ContextData) GiveTraceId(id string) {
	c.TraceId = id
}

func (c *ContextData) GiveUser(user interface{}) {
	c.User = user
}

func (c *ContextData) GiveParams(params interface{}) {
	c.Params = params
}

func (c *ContextData) GiveIP(ip string) {
	c.IP = ip
}

func (c *ContextData) Set(key string, value interface{}) {
	c.Data[key] = value
}

func (c *ContextData) Get(key string, def interface{}) interface{} {
	if value, exists := c.Data[key]; exists {
		return value
	}
	return def
}

type Context interface {
	ILogger
	TakeData() IContextData
	AfterFunc(f func()) (stop func() bool)
	WithTimeout(timeout time.Duration) (ctx Context, cancel context.CancelFunc)
	WithTimeoutCause(timeout time.Duration, cause error) (ctx Context, cancel context.CancelFunc)
	WithCancel() (ctx Context, cancel context.CancelFunc)
	WithCancelCause() (ctx Context, cancel context.CancelCauseFunc)
	WithDeadline(deadline time.Time) (ctx Context, cancel context.CancelFunc)
	WithDeadlineCause(deadline time.Time, cause error) (ctx Context, cancel context.CancelFunc)
}

func NewContext(log ILogger) Context {
	return &defaultContext{
		ctx:     context.Background(),
		ILogger: log,
	}
}

func ctxNewContext(ctx context.Context, log ILogger) Context {
	return &defaultContext{
		ctx:     ctx,
		ILogger: log,
	}
}

type defaultContext struct {
	ILogger
	ctx context.Context
}

func (d defaultContext) TakeData() IContextData {
	return d.ILogger.TakeContextData()
}

func (d defaultContext) AfterFunc(f func()) (stop func() bool) {
	return context.AfterFunc(d.ctx, f)
}

func (d defaultContext) WithTimeout(timeout time.Duration) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithTimeout(d.ctx, timeout)
	return d.makeChildContext(child), cancel
}

func (d defaultContext) WithTimeoutCause(timeout time.Duration, cause error) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithTimeoutCause(d.ctx, timeout, cause)
	return d.makeChildContext(child), cancel
}

func (d defaultContext) WithCancel() (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithCancel(d.ctx)
	return d.makeChildContext(child), cancel
}

func (d defaultContext) WithCancelCause() (ctx Context, cancel context.CancelCauseFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithCancelCause(d.ctx)
	return d.makeChildContext(child), cancel
}

func (d defaultContext) makeChildContext(child context.Context) (ctx Context) {
	return ctxNewContext(child, d.ILogger)
}

func (d defaultContext) WithDeadline(deadline time.Time) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithDeadline(d.ctx, deadline)

	return d.makeChildContext(child), cancel
}

func (d defaultContext) WithDeadlineCause(deadline time.Time, cause error) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithDeadlineCause(d.ctx, deadline, cause)
	return d.makeChildContext(child), cancel
}
