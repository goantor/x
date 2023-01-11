package x

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type IApplications map[string]IApplication

type IApplication interface {
	// Boot 启动
	Boot() error

	// Bootstrap 启动项
	Bootstrap(boots ...IBootstrap) IApplication

	// AutoConfig 自动配置
	//AutoConfig() IApplication

	// Handler 获取路由
	Handler() interface{}

	// Shutdown 关闭服务
	Shutdown()

	// Config 自定义配置
	Config(config interface{}) IApplication
}

func NewHttp(handler http.Handler) IApplication {
	return &DefaultHttp{
		handler: handler,
	}
}

type HttpConfig struct {
	Address      string `json:"address" toml:"address"`
	Port         int    `json:"port" toml:"port"`
	ReadTimeout  int    `json:"read_timeout" toml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" toml:"write_timeout"`
}

func (c HttpConfig) Link() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type DefaultHttp struct {
	server  *http.Server
	handler http.Handler
	config  HttpConfig
}

func (d *DefaultHttp) Bootstrap(boots ...IBootstrap) IApplication {
	for _, boot := range boots {
		boot.Boot(d)
	}

	return d
}

//func (d *DefaultHttp) AutoConfig() IApplication {
//	path := fmt.Sprintf("./config/%s.toml", Env)
//	if _, err := toml.DecodeFile(path, &Configs); err != nil {
//		panic(err)
//	}
//
//	d.config = Configs.Http
//	return d
//}

func (d *DefaultHttp) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := d.server.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func (d *DefaultHttp) Boot() (err error) {
	d.server = &http.Server{
		Addr:         d.config.Link(),
		Handler:      d.handler,
		ReadTimeout:  time.Duration(d.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(d.config.WriteTimeout) * time.Second,
	}

	// 需要修复
	err = d.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return
}

func (d *DefaultHttp) Config(config interface{}) IApplication {
	self, ok := config.(HttpConfig)
	if !ok {
		panic("the config must be a x.HttpConfig")
	}

	d.config = self
	return d
}

func (d *DefaultHttp) Handler() interface{} {
	return d.handler
}
