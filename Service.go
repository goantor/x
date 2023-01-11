package x

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
)

type IBootstrap interface {
	Boot(app IApplication)
}

type IService interface {
	Mounts(name string) IService // 主挂载
	With(names ...string) IService
	Boot()
}

type DefaultService struct {
	main   string
	others []string
}

func (d *DefaultService) Mounts(name string) IService {
	d.main = name
	return d
}

func (d *DefaultService) With(names ...string) IService {
	if d.others == nil {
		d.others = make([]string, 0)
	}

	for _, name := range names {
		d.others = append(d.others, name)
	}

	return d
}

func (d DefaultService) handler(name string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%+v, stack: %s\n", r, debug.Stack())
		}
	}()

	app := Factory(name)
	if app == nil {
		panic(fmt.Sprintf("service {%s} not found", name))
	}
	Shutdown := quit()
	go func() {
		if err := app.Boot(); err != nil {
			panic(fmt.Sprintf("service {%s} boot failed: %v", name, err))
		}
	}()

	<-Shutdown
	app.Shutdown()
}

func (d DefaultService) Boot() {
	// 其他服务先启动
	if len(d.others) > 0 {
		for _, name := range d.others {
			go d.handler(name)
		}
	}

	d.handler(d.main)
}

func quit() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	//signal.Notify(c,os.Signal())
	return c
}

func NewService() IService {
	return &DefaultService{}
}
