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

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{2, 4, 6}

	result := common.Map(input, func(n int) int {
		return n * 2
	})

	require.Equal(t, expected, result)
}

func TestFind(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

	items := []Item{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
	}

	result, found := common.Find(items, func(i Item) bool {
		return i.ID == 2
	})

	require.True(t, found)
	require.Equal(t, "B", result.Name)

	_, found = common.Find(items, func(i Item) bool {
		return i.ID == 999
	})

	require.False(t, found)
}
