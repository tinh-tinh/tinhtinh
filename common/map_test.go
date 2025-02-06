package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

func TestCloneMap(t *testing.T) {
	o := make(map[string]string)
	o["a"] = "vb"

	c := common.CloneMap(o)
	require.Equal(t, o, c)
}

func Test_MergeMaps(t *testing.T) {
	o := make(map[string]string)
	o["a"] = "vb"

	common.MergeMaps(o, map[string]string{"b": "d"})
	require.Equal(t, map[string]string{"a": "vb", "b": "d"}, o)
}
