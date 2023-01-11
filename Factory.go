package x

type ApplicationHandler func() (app IApplication)

var (
	apps   = make(map[string]ApplicationHandler)
	parsed = make(map[string]IApplication)
)

func Factory(name string) IApplication {
	if app, ok := parsed[name]; ok {
		return app
	}

	if app, ok := apps[name]; ok {
		parsed[name] = app()
		return parsed[name]
	}

	return nil
}

func Register(name string, handler func() IApplication) {
	apps[name] = handler
}
