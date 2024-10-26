package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Name(t *testing.T) {
	type Person struct{}
	require.Equal(t, "Person", GetStructName(Person{}))
	require.Equal(t, "Person", GetStructName(&Person{}))
}
