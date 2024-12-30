package x

import (
	"context"
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sync"
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
	//GiveRemind(key string, val interface{})
	GiveUserId(id int)
	GivePlanId(id int)
	GivePhoneMd5(md5 string)
	GiveChannel(channel string)
	TakeActive() string
	GiveActive(active string)
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
		Data:    sync.Map{},
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
	Active    string      `json:"active,omitempty"`
	UserId    int         `json:"user_id,omitempty"`
	PlanId    int         `json:"plan_id,omitempty"`
	PhoneMd5  string      `json:"phone_md5,omitempty"`
	Channel   string      `json:"channel,omitempty"`
	Params    interface{} `json:"params,omitempty"`
	Mark      string      `json:"X_MARK,omitempty"` // 小模块标记 可以清除
	IP        string      `json:"ip,omitempty"`
	Data      sync.Map    `json:"-"`
}

func (c *ContextData) GiveUserId(id int) {
	c.UserId = id
}

func (c *ContextData) GivePlanId(id int) {
	c.PlanId = id
}

func (c *ContextData) GivePhoneMd5(md5 string) {
	c.PhoneMd5 = md5
}

func (c *ContextData) GiveChannel(channel string) {
	c.Channel = channel
}

func (c *ContextData) GiveActive(active string) {
	c.Active = active
}

func (c *ContextData) TakeActive() string {
	return c.Active
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

//func (c *ContextData) GiveRemind(key string, val interface{}) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	c.Remind[key] = val
//}

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
	c.Data.Store(key, value)
}

func (c *ContextData) Get(key string, def interface{}) interface{} {
	if value, exists := c.Data.Load(key); exists {
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

	// GiveRemoteRequestTimeout 设定远程请求超时时间
	GiveRemoteRequestTimeout(timeout time.Duration)

	// TakeRemoteRequestTimeout 获取远程请求超时时间
	TakeRemoteRequestTimeout(def time.Duration) time.Duration

	RecordTime()

	TakeUsedTime() time.Duration
}

func NewContext(log ILogger) Context {
	return &defaultContext{
		Context:  context.Background(),
		ILogger:  log,
		UseTimer: NewUseTimer(),
	}
}

func NewContextWithLog(log *logrus.Logger) Context {
	return NewContext(
		NewLogger(log),
	)
}

func ctxNewContext(ctx context.Context, log ILogger) Context {
	return &defaultContext{
		Context:  ctx,
		ILogger:  log,
		UseTimer: NewUseTimer(),
	}
}

type defaultContext struct {
	ILogger
	ILocker
	context.Context

	RemoteRequestTimeout time.Duration `json:"remote_request_timeout"` //
	UseTimer             *UseTimer
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

func (d *defaultContext) GiveRemoteRequestTimeout(timeout time.Duration) {
	d.RemoteRequestTimeout = timeout
}

func (d *defaultContext) TakeRemoteRequestTimeout(def time.Duration) time.Duration {
	if d.RemoteRequestTimeout == 0 {
		return def
	}
	return d.RemoteRequestTimeout
}

func (d *defaultContext) RecordTime() {
	d.UseTimer.Record()
}

func (d *defaultContext) TakeUsedTime() time.Duration {
	return d.UseTimer.TakeUsed()
}
