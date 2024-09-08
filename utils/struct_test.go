package utils

import "testing"

type User struct{}
type Post struct{}
type TinhTinh struct{}

func Test_GetNameStruct(t *testing.T) {
	var testcases = []struct {
		input interface{}
		want  string
	}{
		{input: &User{}, want: "User"},
		{input: &Post{}, want: "Post"},
		{input: &TinhTinh{}, want: "TinhTinh"},
	}

	for _, tc := range testcases {
		t.Run("test case", func(t *testing.T) {
			got := GetNameStruct(tc.input)
			if got != tc.want {
				t.Errorf("expect %s, but got %s", tc.want, got)
			}
		})
	}
}
