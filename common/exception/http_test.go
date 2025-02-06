package exception_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Exception(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("bad-request", func(ctx core.Ctx) error {
			panic(exception.BadRequest("bad request"))
		})

		ctrl.Get("unauthorized", func(ctx core.Ctx) error {
			panic(exception.Unauthorized("unauthorized"))
		})

		ctrl.Get("forbidden", func(ctx core.Ctx) error {
			panic(exception.Forbidden("forbidden"))
		})

		ctrl.Get("not-found", func(ctx core.Ctx) error {
			panic(exception.NotFound("not found"))
		})

		ctrl.Get("method-not-allowed", func(ctx core.Ctx) error {
			panic(exception.MethodNotAllowed("method not allowed"))
		})

		ctrl.Get("not-acceptable", func(ctx core.Ctx) error {
			panic(exception.NotAcceptable("not acceptable"))
		})

		ctrl.Get("request-timeout", func(ctx core.Ctx) error {
			panic(exception.RequestTimeout("request timeout"))
		})

		ctrl.Get("conflict", func(ctx core.Ctx) error {
			panic(exception.Conflict("conflict"))
		})

		ctrl.Get("gone", func(ctx core.Ctx) error {
			panic(exception.Gone("gone"))
		})

		ctrl.Get("length-required", func(ctx core.Ctx) error {
			panic(exception.LengthRequired("length required"))
		})

		ctrl.Get("precondition-failed", func(ctx core.Ctx) error {
			panic(exception.PreconditionFailed("precondition failed"))
		})

		ctrl.Get("payload-too-large", func(ctx core.Ctx) error {
			panic(exception.ContentTooLarge("payload too large"))
		})

		ctrl.Get("uri-too-long", func(ctx core.Ctx) error {
			panic(exception.RequestURITooLong("uri too long"))
		})

		ctrl.Get("unsupported-media-type", func(ctx core.Ctx) error {
			panic(exception.UnsupportedMediaType("unsupported media type"))
		})

		ctrl.Get("range-not-satisfiable", func(ctx core.Ctx) error {
			panic(exception.RequestedRangeNotSatisfiable("range not satisfiable"))
		})

		ctrl.Get("expectation-failed", func(ctx core.Ctx) error {
			panic(exception.ExpectationFailed("expectation failed"))
		})

		ctrl.Get("im-a-teapot", func(ctx core.Ctx) error {
			panic(exception.IamATeapot("im a teapot"))
		})

		ctrl.Get("unprocessable-entity", func(ctx core.Ctx) error {
			panic(exception.UnprocessableEntity("unprocessable entity"))
		})

		ctrl.Get("locked", func(ctx core.Ctx) error {
			panic(exception.Locked("locked"))
		})

		ctrl.Get("failed-dependency", func(ctx core.Ctx) error {
			panic(exception.FailedDependency("failed dependency"))
		})

		ctrl.Get("upgrade-required", func(ctx core.Ctx) error {
			panic(exception.UpgradeRequired("upgrade required"))
		})

		ctrl.Get("precondition-required", func(ctx core.Ctx) error {
			panic(exception.PreconditionRequired("precondition required"))
		})

		ctrl.Get("too-many-requests", func(ctx core.Ctx) error {
			panic(exception.TooManyRequests("too many requests"))
		})

		ctrl.Get("internal-server-error", func(ctx core.Ctx) error {
			panic(exception.InternalServer("internal server error"))
		})

		ctrl.Get("not-implemented", func(ctx core.Ctx) error {
			panic(exception.NotImplemented("not implemented"))
		})

		ctrl.Get("bad-gateway", func(ctx core.Ctx) error {
			panic(exception.BadGateway("bad gateway"))
		})

		ctrl.Get("service-unavailable", func(ctx core.Ctx) error {
			panic(exception.ServiceUnavailable("service unavailable"))
		})

		ctrl.Get("gateway-timeout", func(ctx core.Ctx) error {
			panic(exception.GatewayTimeout("gateway timeout"))
		})

		ctrl.Get("http-version-not-supported", func(ctx core.Ctx) error {
			panic(exception.HttpVersionNotSupported("http version not supported"))
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/bad-request")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unauthorized")
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/forbidden")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/not-found")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/method-not-allowed")
	require.Nil(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/not-acceptable")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotAcceptable, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/request-timeout")
	require.Nil(t, err)
	require.Equal(t, http.StatusRequestTimeout, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/conflict")
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/gone")
	require.Nil(t, err)
	require.Equal(t, http.StatusGone, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/length-required")
	require.Nil(t, err)
	require.Equal(t, http.StatusLengthRequired, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/precondition-failed")
	require.Nil(t, err)
	require.Equal(t, http.StatusPreconditionFailed, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/payload-too-large")
	require.Nil(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/uri-too-long")
	require.Nil(t, err)
	require.Equal(t, http.StatusRequestURITooLong, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unsupported-media-type")
	require.Nil(t, err)
	require.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/range-not-satisfiable")
	require.Nil(t, err)
	require.Equal(t, http.StatusRequestedRangeNotSatisfiable, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/expectation-failed")
	require.Nil(t, err)
	require.Equal(t, http.StatusExpectationFailed, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/im-a-teapot")
	require.Nil(t, err)
	require.Equal(t, http.StatusTeapot, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unprocessable-entity")
	require.Nil(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/locked")
	require.Nil(t, err)
	require.Equal(t, http.StatusLocked, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/failed-dependency")
	require.Nil(t, err)
	require.Equal(t, http.StatusFailedDependency, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/upgrade-required")
	require.Nil(t, err)
	require.Equal(t, http.StatusUpgradeRequired, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/precondition-required")
	require.Nil(t, err)
	require.Equal(t, http.StatusPreconditionRequired, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/too-many-requests")
	require.Nil(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/internal-server-error")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/not-implemented")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/bad-gateway")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/service-unavailable")
	require.Nil(t, err)
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/gateway-timeout")
	require.Nil(t, err)
	require.Equal(t, http.StatusGatewayTimeout, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/http-version-not-supported")
	require.Nil(t, err)
	require.Equal(t, http.StatusHTTPVersionNotSupported, resp.StatusCode)
}
