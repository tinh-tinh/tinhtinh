package cookie_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/cookie"
)

func Test_Encode(t *testing.T) {
	require.Equal(t, "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92amY=", cookie.Encode([]byte("oievrervierjvoiejrojojovjf")))
	decode, err := cookie.Decode("b2lldnJlcnZpZXJqdm9pZWpyb2pvam92amY=")
	require.Nil(t, err)
	require.Equal(t, "oievrervierjvoiejrojojovjf", string(decode))

	_, err = cookie.Decode("Tôi Tích Ta Tu Tiên")
	require.NotNil(t, err)

	sCookie := &cookie.SecureCookie{Key: "add"}
	_, err = sCookie.Encrypt("abc")
	require.NotNil(t, err)

	_, err = sCookie.Decrypt("avv")
	require.NotNil(t, err)

	sCookie2 := &cookie.SecureCookie{Key: "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92amY="}
	_, err = sCookie2.Decrypt("Tôi Tích Ta Tu Tiên")
	require.NotNil(t, err)
}
