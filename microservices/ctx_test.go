package microservices_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func Test_Ctx(t *testing.T) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	svc := tcp.New(func() core.Module { return nil }, tcp.Options{Addr: "127.0.0.1:5173"})
	input := `"email": "xyz@gmail.com", "password": "12345678@Tc"`

	// Wrap the string into valid JSON format
	jsonString := fmt.Sprintf("{%s}", input)

	message := microservices.Message{
		Type:  microservices.RPC,
		Event: "test",
		Headers: map[string]string{
			"key": "value",
		},
		Data: jsonString,
	}
	ctx := microservices.NewCtx(message, svc)

	payload := ctx.Payload(&User{}).(*User)
	fmt.Println(payload)
	require.Equal(t, "xyz@gmail.com", payload.Email)
	require.Equal(t, "12345678@Tc", payload.Password)

	ctx.Set("key", "value")
	require.Equal(t, "value", ctx.Get("key"))

	headers := ctx.Headers("key")
	require.Equal(t, "value", headers)
}
