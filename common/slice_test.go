package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

func Test_Filter(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	res := common.Filter(data, func(item int) bool {
		return item%2 == 0
	})

	require.Equal(t, []int{2, 4}, res)
}

func Test_Remove(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	res := common.Remove(data, func(item int) bool {
		return item%2 == 0
	})

	require.Equal(t, []int{1, 3, 5}, res)
}
