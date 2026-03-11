package core_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// helper: a simple multi-module app used across visualize tests.
func buildVisualizeApp() *core.App {
	const UserSvc core.Provide = "UserService"
	const AuthSvc core.Provide = "AuthService"

	userProvider := func(module core.Module) core.Provider {
		return module.NewProvider(core.ProviderOptions{
			Name:  UserSvc,
			Value: "user-value",
		})
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("users")
		ctrl.Get("", func(ctx core.Ctx) error { return ctx.JSON(core.Map{"data": "list"}) })
		ctrl.Post("", func(ctx core.Ctx) error { return ctx.JSON(core.Map{"data": "created"}) })
		ctrl.Get(":id", func(ctx core.Ctx) error { return ctx.JSON(core.Map{"data": "one"}) })
		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{userController},
			Providers:   []core.Providers{userProvider},
		})
		mod.Export(UserSvc)
		return mod
	}

	authProvider := func(module core.Module) core.Provider {
		return module.NewProvider(core.ProviderOptions{
			Name:  AuthSvc,
			Value: "auth-value",
		})
	}

	authController := func(module core.Module) core.Controller {
		ctrl := module.NewController("auth")
		ctrl.Post("login", func(ctx core.Ctx) error { return ctx.JSON(core.Map{"data": "token"}) })
		ctrl.Post("logout", func(ctx core.Ctx) error { return ctx.JSON(core.Map{"data": "ok"}) })
		return ctrl
	}

	authModule := func(module core.Module) core.Module {
		return module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{authController},
			Providers:   []core.Providers{authProvider},
		})
	}

	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{userModule, authModule},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")
	return app
}

func Test_Visualize_GetTree(t *testing.T) {
	app := buildVisualizeApp()

	tree := app.GetTree()
	require.NotNil(t, tree)
	require.Equal(t, "AppModule", tree.Name)
	require.Len(t, tree.Imports, 2)

	// First import is userModule
	userNode := tree.Imports[0]
	require.NotEmpty(t, userNode.Name)
	require.NotEmpty(t, userNode.Controllers)
	require.Equal(t, "users", userNode.Controllers[0].Name)
	require.Len(t, userNode.Controllers[0].Routes, 3)

	// Second import is authModule
	authNode := tree.Imports[1]
	require.NotEmpty(t, authNode.Name)
	require.Equal(t, "auth", authNode.Controllers[0].Name)
	require.Len(t, authNode.Controllers[0].Routes, 2)
}

func Test_Visualize_HTML(t *testing.T) {
	app := buildVisualizeApp()

	outputPath := "/tmp/tinhtinh_tree_test.html"
	t.Cleanup(func() { os.Remove(outputPath) })

	err := app.Visualize(outputPath)
	require.NoError(t, err)

	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	require.Greater(t, info.Size(), int64(1000), "HTML file should be non-trivial")

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "TinhTinh Module Tree")
	require.Contains(t, string(content), "AppModule")
	require.Contains(t, string(content), "users")
	require.Contains(t, string(content), "auth")
}

// Test_Visualize_GetTree_NonDynamicModule verifies that GetTree returns nil
// when the App's module has been replaced with a plain (non-DynamicModule) value.
func Test_Visualize_GetTree_NonDynamicModule(t *testing.T) {
	// Build a valid app first, then swap Module with a non-DynamicModule recorder.
	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{})
	}
	app := core.CreateFactory(appModule)

	// Replace the module with nil – nil satisfies the Module interface but the
	// type assertion to *DynamicModule inside GetTree will fail, so nil is returned.
	app.Module = nil

	tree := app.GetTree()
	require.Nil(t, tree, "GetTree must return nil for a non-DynamicModule module")
}

// Test_Visualize_NilTree verifies that Visualize returns nil (no error) and
// does NOT create any file when GetTree returns nil.
func Test_Visualize_NilTree(t *testing.T) {
	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{})
	}
	app := core.CreateFactory(appModule)
	app.Module = nil // force GetTree → nil

	outputPath := "/tmp/tinhtinh_nil_tree_test.html"
	t.Cleanup(func() { os.Remove(outputPath) })

	err := app.Visualize(outputPath)
	require.NoError(t, err, "Visualize should succeed (no-op) when tree is nil")

	_, statErr := os.Stat(outputPath)
	require.True(t, os.IsNotExist(statErr), "no HTML file should be written when tree is nil")
}

// Test_Visualize_UnwritablePath verifies that Visualize returns an error when
// the output path cannot be created (e.g. directory does not exist).
func Test_Visualize_UnwritablePath(t *testing.T) {
	app := buildVisualizeApp()

	// Use a path whose parent directory does not exist.
	unwritable := "/tmp/nonexistent_dir_tinhtinh/tree.html"

	err := app.Visualize(unwritable)
	require.Error(t, err, "Visualize must return an error for an unwritable path")
}

// Test_Visualize_EmptyModule verifies that GetTree and Visualize work correctly
// for an app whose root module has no controllers, providers, or imports.
func Test_Visualize_EmptyModule(t *testing.T) {
	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{})
	}
	app := core.CreateFactory(appModule)

	tree := app.GetTree()
	require.NotNil(t, tree)
	require.Empty(t, tree.Controllers, "empty module should have no controllers")
	require.Empty(t, tree.Providers, "empty module should have no providers")
	require.Empty(t, tree.Imports, "empty module should have no imports")

	outputPath := "/tmp/tinhtinh_empty_module_test.html"
	t.Cleanup(func() { os.Remove(outputPath) })

	err := app.Visualize(outputPath)
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "TinhTinh Module Tree")
}
