package logger

import "github.com/tinh-tinh/tinhtinh/v2/core"

const LOGGER core.Provide = "LOGGER"

// Module creates a new logger module, which is a global module.
//
// It takes an Options struct as a parameter, which is used to create the logger.
// The logger is created with the given options and is registered as a provider in the module
// with the name LOGGER. The logger is also exported by the module.
func Module(opt Options) core.Modules {
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

// InjectLog injects the logger from the module to the caller. It returns a
// pointer to the logger, or nil if the logger is not found.
func InjectLog(module *core.DynamicModule) *Logger {
	return module.Ref(LOGGER).(*Logger)
}
