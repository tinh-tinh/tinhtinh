package logger

import "github.com/tinh-tinh/tinhtinh/core"

const LOGGER core.Provide = "LOGGER"

func Module(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		loggerModule := module.New(core.NewModuleOptions{
			Scope: core.Global,
		})

		loggerModule.NewProvider(core.ProviderOptions{
			Name:  LOGGER,
			Value: Create(opt),
		})

		loggerModule.Export(LOGGER)
		return loggerModule
	}
}

func InjectLog(module *core.DynamicModule) *Logger {
	return module.Ref(LOGGER).(*Logger)
}
