package tcp_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

type User struct {
	Email    string `json:"email" validate:"required,isEmail"`
	Password string `json:"password" validate:"required,isStrongPassword"`
}

func AuthApp(addr string) *core.App {
	const USER_SERVICE core.Provide = "users_service"
	type UserService struct {
		users []*User
		total int64
	}

	service := func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: USER_SERVICE,
			Value: &UserService{
				users: []*User{},
				total: 0,
			},
		})

		return prd
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("auth")
		client := microservices.InjectClient(module, microservices.TCP)

		ctrl.Post("register", func(ctx core.Ctx) error {
			userService := module.Ref(USER_SERVICE).(*UserService)
			user := &User{}
			if err := ctx.BodyParser(user); err != nil {
				return ctx.Status(400).JSON(core.Map{"error": err.Error()})
			}
			tenantID := ctx.Headers("x-tenant-id")

			userService.users = append(userService.users, user)
			atomic.AddInt64(&userService.total, 1)

			go client.Send("user.created", user, microservices.Header{"x-tenant-id": tenantID})
			return ctx.JSON(core.Map{"data": user})
		})

		ctrl.Post("login", func(ctx core.Ctx) error {
			userService := module.Ref(USER_SERVICE).(*UserService)
			user := &User{}
			if err := ctx.BodyParser(user); err != nil {
				return ctx.Status(400).JSON(core.Map{"error": err.Error()})
			}

			for _, u := range userService.users {
				if u.Email == user.Email && u.Password == user.Password {
					tenantID := ctx.Headers("x-tenant-id")

					go client.Publish("user.*", user, microservices.Header{"x-tenant-id": tenantID})
					return ctx.JSON(core.Map{"data": u})
				}
			}

			return ctx.Status(400).JSON(core.Map{"error": "user not found"})
		})

		ctrl.Get("total", func(ctx core.Ctx) error {
			userService := module.Ref(USER_SERVICE).(*UserService)
			return ctx.JSON(core.Map{"total": userService.total})
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: microservices.TCP,
					Transport: tcp.NewClient(tcp.Options{
						Addr:    addr,
						Timeout: 200 * time.Millisecond,
						Config: microservices.Config{
							CompressAlg: compress.Gzip,
							RetryOptions: microservices.RetryOptions{
								Retry: 3,
								Delay: 100 * time.Millisecond,
							},
						},
					}),
				}),
			},
			Controllers: []core.Controllers{controller},
			Providers: []core.Providers{
				service,
			},
		})

		module.NewController("auth")
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("auth-api")

	return app
}

func DirectoryApp() *core.App {
	const DIRECTORY_SERVICE core.Provide = "directory_service"
	type DirectoryService struct {
		directories     map[string][]*User
		currentlyLogged map[string]*User
	}

	middleware := func(ctx microservices.Ctx) error {
		tenantID := ctx.Headers("x-tenant-id")
		if tenantID != "" {
			ctx.Set("tenant", tenantID)
		}
		return ctx.Next()
	}

	service := func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: DIRECTORY_SERVICE,
			Value: &DirectoryService{
				directories:     make(map[string][]*User),
				currentlyLogged: make(map[string]*User),
			},
		})

		return prd
	}

	handlers := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{}).Use(middleware).Registry()

		directoryService := module.Ref(DIRECTORY_SERVICE).(*DirectoryService)
		handler.OnResponse("user.created", func(ctx microservices.Ctx) error {
			tenantID := ctx.Get("tenant").(string)

			if directoryService.directories[tenantID] == nil {
				directoryService.directories[tenantID] = []*User{}
			}

			var payload *User
			err := ctx.PayloadParser(&payload)
			if err != nil {
				fmt.Println(err)
				return err
			}
			directoryService.directories[tenantID] = append(directoryService.directories[tenantID], payload)

			fmt.Printf("Receive payload is %v\n", payload)
			return nil
		})

		handler.OnEvent("user.logined", func(ctx microservices.Ctx) error {
			tenantID := ctx.Get("tenant").(string)

			fmt.Println("User Logged Data:", ctx.Payload(), tenantID)
			var payload *User
			err := ctx.PayloadParser(&payload)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if payload != nil {
				directoryService.currentlyLogged[tenantID] = payload
			}
			return nil
		})

		return handler
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("directory")

		ctrl.Get("total", func(ctx core.Ctx) error {
			directoryService := module.Ref(DIRECTORY_SERVICE).(*DirectoryService)
			tenant := ctx.Headers("x-tenant-id")
			return ctx.JSON(core.Map{"data": len(directoryService.directories[tenant])})
		})

		ctrl.Get("current", func(ctx core.Ctx) error {
			directoryService := module.Ref(DIRECTORY_SERVICE).(*DirectoryService)

			tenant := ctx.Headers("x-tenant-id")
			return ctx.JSON(core.Map{"data": directoryService.currentlyLogged[tenant]})
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:     []core.Modules{microservices.Register()},
			Controllers: []core.Controllers{controller},
			Providers:   []core.Providers{service, handlers},
		})
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("directory-api")

	return app
}

func Test_Tenant(t *testing.T) {
	directoryApp := DirectoryApp()
	directoryApp.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:4001",
		Config: microservices.Config{
			CompressAlg: compress.Gzip,
		},
	}))
	directoryApp.StartAllMicroservices()

	testServerDirectory := httptest.NewServer(directoryApp.PrepareBeforeListen())
	defer testServerDirectory.Close()

	time.Sleep(100 * time.Millisecond)

	authApp := AuthApp("localhost:4001")
	testServerAuth := httptest.NewServer(authApp.PrepareBeforeListen())
	defer testServerAuth.Close()

	testClientAuth := testServerAuth.Client()
	testClientDirectory := testServerDirectory.Client()

	req, err := http.NewRequest("POST", testServerAuth.URL+"/auth-api/auth/register", strings.NewReader(`{"email": "xyz@gmail.com", "password": "12345678@Tc"}`))
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "tenant1")

	resp, err := testClientAuth.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(100 * time.Millisecond)

	req, err = http.NewRequest("GET", testServerDirectory.URL+"/directory-api/directory/total", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "tenant1")
	resp, err = testClientDirectory.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":1}`, string(data))

	req, err = http.NewRequest("POST", testServerAuth.URL+"/auth-api/auth/login", strings.NewReader(`{"email": "xyz@gmail.com", "password": "12345678@Tc"}`))
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "tenant1")
	resp, err = testClientAuth.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(100 * time.Millisecond)

	req, err = http.NewRequest("GET", testServerDirectory.URL+"/directory-api/directory/current", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "tenant1")
	resp, err = testClientDirectory.Do(req)
	require.Nil(t, err)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{"email":"xyz@gmail.com","password":"12345678@Tc"}}`, string(data))
}
