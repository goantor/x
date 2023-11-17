package x

/*
	x包定义框架基础能力,也可以配合使用
	DefaultConfigs 配置
	Context 中间上下文
	Error 错误
	defaultLog 日志配置
	Validate 验证
*/

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goantor/logs"
	"time"
)

type Roboter interface {
	Send(message string) error
}

type Context interface {
	logs.Logger
	Set(key string, value interface{})
	Get(key string, def interface{}) interface{}
	Response(code int, h H)
}

type GinContext struct {
	logs.Logger
	gtx *gin.Context
	ctx context.Context
}

func (g GinContext) Set(key string, value interface{}) {
	g.gtx.Set(key, value)
}

func (g GinContext) Get(key string, def interface{}) interface{} {
	if value, exists := g.gtx.Get(key); exists {
		return value
	}

	return def
}

func (g GinContext) Context() context.Context {
	return g.ctx
}

func (g GinContext) Timeout(duration time.Duration) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(g.ctx, duration)
}

func (g GinContext) Response(code int, h H) {
	if g.gtx.Writer.Written() {
		return
	}

	g.gtx.AbortWithStatusJSON(code, h)
}

func NewContextWithGin(ctx *gin.Context, log logs.Logger) Context {
	return &GinContext{gtx: ctx, Logger: log, ctx: context.Background()}
}
