package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct{}
type Post struct{}
type TinhTinh struct{}

func Test_GetNameStruct(t *testing.T) {
	t.Parallel()

	require.Equal(t, "User", GetNameStruct(&User{}))
	require.Equal(t, "Post", GetNameStruct(&Post{}))
	require.Equal(t, "TinhTinh", GetNameStruct(&TinhTinh{}))
}
