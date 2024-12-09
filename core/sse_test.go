package core_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func TestSseHandler(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("")

		ctrl.Sse("events", func(broker *core.SseBroker) {
			for i := 0; i < 3; i++ {
				broker.Messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

				// Print a nice log message and sleep for 5s.
				log.Printf("Sent message %d ", i)
				time.Sleep(1e9)
			}

			broker.Close()
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)

	req := httptest.NewRequest("GET", "/events/", nil)
	rr := httptest.NewRecorder()

	handler := app.PrepareBeforeListen()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "text/event-stream", rr.Header().Get("Content-Type"))
}
