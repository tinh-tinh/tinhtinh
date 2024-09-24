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

// OnInit registers the given hooks to be run when the module is initialized.
// The hooks are run in the order they are added to the module.
func (m *DynamicModule) OnInit(hooks ...HookModule) *DynamicModule {
	m.hooks = append(m.hooks, hooks...)
	return m
}

// init initializes the module by running the OnInit hooks in the order they were added to the module.
func (m *DynamicModule) init() {
	for _, v := range m.hooks {
		v(m)
	}
}

// BeforeShutdown registers the given hooks to be run before the server is shut down.
// The hooks are run in the order they are added to the App instance.
func (app *App) BeforeShutdown(hooks ...HookFnc) *App {
	for _, hook := range hooks {
		app.hooks = append(app.hooks, &Hook{
			fnc:   hook,
			RunAt: BEFORE_SHUTDOWN,
		})
	}
	return app
}

// AfterShutdown registers the given hooks to be run after the server is shut down.
// The hooks are run in the order they are added to the App instance.
func (app *App) AfterShutdown(hooks ...HookFnc) *App {
	for _, hook := range hooks {
		app.hooks = append(app.hooks, &Hook{
			fnc:   hook,
			RunAt: AFTER_SHUTDOWN,
		})
	}
	return app
}
