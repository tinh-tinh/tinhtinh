package tcp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func Test_Timeout(t *testing.T) {
	directoryApp := DirectoryApp()
	directoryApp.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:4002",
	}))
	directoryApp.StartAllMicroservices()

	testServerDirectory1 := httptest.NewServer(directoryApp.PrepareBeforeListen())

	time.Sleep(100 * time.Millisecond)

	authApp := AuthApp("localhost:4002")
	testServerAuth := httptest.NewServer(authApp.PrepareBeforeListen())
	defer testServerAuth.Close()

	testClientAuth := testServerAuth.Client()
	req, err := http.NewRequest("POST", testServerAuth.URL+"/auth-api/auth/register", strings.NewReader(`{"email": "xyz@gmail.com", "password": "12345678@Tc"}`))
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "tenant")

	resp, err := testClientAuth.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	testServerDirectory1.Close()
	time.Sleep(1000 * time.Millisecond)

	req, err = http.NewRequest("POST", testServerAuth.URL+"/auth-api/auth/register", strings.NewReader(`{"email": "abc@gmail.com", "password": "12345678@Tc"}`))
	require.Nil(t, err)
	req.Header.Set("x-tenant-id", "tenant")

	resp, err = testClientAuth.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(200 * time.Millisecond)
}
