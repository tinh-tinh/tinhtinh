package microservices_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/microservices/tcp"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
)

func Test_Ctx(t *testing.T) {
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	svc := tcp.NewServer(tcp.Options{Addr: "127.0.0.1:5173"})
	input := `"email": "xyz@gmail.com", "password": "12345678@Tc"`

	// Wrap the string into valid JSON format
	jsonString := fmt.Sprintf("{%s}", input)

	message := microservices.Message{
		Event: "test",
		Headers: map[string]string{
			"key": "value",
		},
		Data: jsonString,
	}
	ctx := microservices.NewCtx(message, svc)

	var payload User
	err := ctx.PayloadParser(&payload)
	require.Nil(t, err)
	require.Equal(t, "xyz@gmail.com", payload.Email)
	require.Equal(t, "12345678@Tc", payload.Password)

	ctx.Set("key", "value")
	require.Equal(t, "value", ctx.Get("key"))

	headers := ctx.Headers("key")
	require.Equal(t, "value", headers)

	message2 := microservices.Message{
		Event: "test",
		Headers: map[string]string{
			"key": "value",
		},
		Data: &User{
			Email:    "abc@gmail.com",
			Password: "12345678",
		},
	}
	encoder, err := compress.Encode(message2.Data, compress.Gzip)
	require.Nil(t, err)

	decoder, err := compress.Decode(encoder, compress.Gzip)
	require.Nil(t, err)

	message2.Bytes = decoder
	message2.Data = nil

	ctx2 := microservices.NewCtx(message2, svc)
	var payload2 User
	err = ctx2.PayloadParser(&payload2)
	require.Nil(t, err)
	require.Equal(t, "xyz@gmail.com", payload.Email)
	require.Equal(t, "12345678@Tc", payload.Password)
}
