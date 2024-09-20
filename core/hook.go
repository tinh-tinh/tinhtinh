package core

type HookFnc func()

type RunAt string

const (
	BEFORE_SHUTDOWN RunAt = "BeforeShutdown"
	AFTER_SHUTDOWN  RunAt = "AfterShutdown"
)

type Hook struct {
	fnc   HookFnc
	RunAt RunAt
}

type HookModule func(module *DynamicModule)

func (m *DynamicModule) OnInit(hooks ...HookModule) *DynamicModule {
	m.hooks = append(m.hooks, hooks...)
	return m
}

func (m *DynamicModule) init() {
	for _, v := range m.hooks {
		v(m)
	}
}

func (app *App) BeforeShutdown(hooks ...HookFnc) *App {
	for _, hook := range hooks {
		app.hooks = append(app.hooks, &Hook{
			fnc:   hook,
			RunAt: BEFORE_SHUTDOWN,
		})
	}
	return app
}

func (app *App) AfterShutdown(hooks ...HookFnc) *App {
	for _, hook := range hooks {
		app.hooks = append(app.hooks, &Hook{
			fnc:   hook,
			RunAt: AFTER_SHUTDOWN,
		})
	}
	return app
}
