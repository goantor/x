package x

import (
	"context"
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type ILocker interface {
	Locked(mark string, duration time.Duration) bool
	UnLock(mark string) error
}

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
	GiveRemind(key string, val interface{})
	GiveParams(params interface{})
	GiveMaskParams(params interface{})
	GiveIP(ip string)
	Set(key string, value interface{})
	Get(key string, def interface{}) interface{}
	GiveMark(mark string)
	TakeMark() string
	CleanMark()
}

func makeTraceId() string {
	buf := make([]byte, 32)
	u := uuid.NewV4().Bytes()
	hex.Encode(buf, u)
	return string(buf)
}

func NewContextData() IContextData {
	return &ContextData{
		Remind:  make(H),
		Data:    make(H),
		TraceId: makeTraceId(),
	}
}

type ContextData struct {
	//UseMasker   bool		`json:UseMask`
	Service   string      `json:"service,omitempty"`
	Module    string      `json:"module,omitempty"`
	Action    string      `json:"action,omitempty"`
	TraceId   string      `json:"trace_id,omitempty"`
	RequestId string      `json:"request_id,omitempty"`
	Remind    H           `json:"X_REMIND,omitempty"` // 这个始终标记不可以清除, 并且为大写，防止与其他数据冲突
	Params    interface{} `json:"params,omitempty"`
	Mark      string      `json:"X_MARK,omitempty"` // 小模块标记 可以清除
	IP        string      `json:"ip,omitempty"`
	Data      H           `json:"-"`
}

func (c *ContextData) GiveMark(mark string) {
	c.Mark = fmt.Sprintf("__%s__", mark)
}

func (c *ContextData) TakeMark() string {
	return c.Mark
}

func (c *ContextData) CleanMark() {
	c.Mark = ""
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

func (c *ContextData) GiveRemind(key string, val interface{}) {
	if _, exists := c.Remind[key]; !exists {
		c.Remind[key] = val
	}
}

func (c *ContextData) GiveParams(params interface{}) {
	c.Params = params
}

func (c *ContextData) GiveMaskParams(params interface{}) {
	masker := MaskReflect{}
	c.Params = masker.MakeMask(params)
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
	ILocker
	WithLocker(locker ILocker)
	TakeContext() context.Context
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
		Context: context.Background(),
		ILogger: log,
	}
}

func NewContextWithLog(log *logrus.Logger) Context {
	return NewContext(
		NewLogger(log),
	)
}

func ctxNewContext(ctx context.Context, log ILogger) Context {
	return &defaultContext{
		Context: ctx,
		ILogger: log,
	}
}

type defaultContext struct {
	ILogger
	ILocker
	context.Context
}

func (d *defaultContext) WithLocker(locker ILocker) {
	d.ILocker = locker
}

func (d *defaultContext) TakeContext() context.Context {
	return d.Context
}

func (d *defaultContext) TakeData() IContextData {
	return d.ILogger.TakeContextData()
}

func (d *defaultContext) AfterFunc(f func()) (stop func() bool) {
	return context.AfterFunc(d.Context, f)
}

func (d *defaultContext) WithTimeout(timeout time.Duration) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithTimeout(d.Context, timeout)
	return d.makeChildContext(child), cancel
}

func (d *defaultContext) WithTimeoutCause(timeout time.Duration, cause error) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithTimeoutCause(d.Context, timeout, cause)
	return d.makeChildContext(child), cancel
}

func (d *defaultContext) WithCancel() (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithCancel(d.Context)
	return d.makeChildContext(child), cancel
}

func (d *defaultContext) WithCancelCause() (ctx Context, cancel context.CancelCauseFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithCancelCause(d.Context)
	return d.makeChildContext(child), cancel
}

func (d *defaultContext) makeChildContext(child context.Context) (ctx Context) {
	return ctxNewContext(child, d.ILogger)
}

func (d *defaultContext) WithDeadline(deadline time.Time) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithDeadline(d.Context, deadline)

	return d.makeChildContext(child), cancel
}

func (d *defaultContext) WithDeadlineCause(deadline time.Time, cause error) (ctx Context, cancel context.CancelFunc) {
	var (
		child context.Context
	)

	child, cancel = context.WithDeadlineCause(d.Context, deadline, cause)
	return d.makeChildContext(child), cancel
}
