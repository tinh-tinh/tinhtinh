package validator

import (
	"fmt"
	"testing"
)

func Test_IsAlpha(t *testing.T) {
	testcase := []struct {
		str  string
		want bool
	}{
		{
			str:  "test",
			want: true,
		},
		{
			str:  "1233",
			want: false,
		},
		{
			str:  "test123!",
			want: false,
		},
		{
			str:  "test123!@#",
			want: false,
		},
	}

	for _, tc := range testcase {
		name := fmt.Sprintf("Test Case %s", tc.str)
		t.Run(name, func(t *testing.T) {
			got := IsAlpha(tc.str)
			if got != tc.want {
				t.Errorf("expect %t, but got %t", tc.want, got)
			}
		})
	}
}

func Test_IsAlphanumeric(t *testing.T) {
	testcases := []struct {
		str  string
		want bool
	}{
		{
			str:  "test",
			want: true,
		},
		{
			str:  "1233",
			want: true,
		},
		{
			str:  "test123!",
			want: false,
		},
		{
			str:  "test123!@#",
			want: false,
		},
		{
			str:  "test123!@#123",
			want: false,
		},
	}

	for _, tc := range testcases {
		name := fmt.Sprintf("Test Case %s", tc.str)
		t.Run(name, func(t *testing.T) {
			got := IsAlphanumeric(tc.str)
			if got != tc.want {
				t.Errorf("expect %t, but got %t", tc.want, got)
			}
		})
	}
}
