package cookie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Encode(t *testing.T) {
	require.Equal(t, "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92amY=", Encode([]byte("oievrervierjvoiejrojojovjf")))
	decode, err := Decode("b2lldnJlcnZpZXJqdm9pZWpyb2pvam92amY=")
	require.Nil(t, err)
	require.Equal(t, "oievrervierjvoiejrojojovjf", string(decode))

	_, err = Decode("Tôi Tích Ta Tu Tiên")
	require.NotNil(t, err)
}
