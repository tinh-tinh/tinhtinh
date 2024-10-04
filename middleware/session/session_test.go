package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Session(t *testing.T) {
	s := New(Options{
		GeneratorID: func() string {
			return time.Now().String()
		},
	})

	cookie := s.Set("abc", "mno")
	require.Equal(t, "mno", s.Get(cookie.Value))
}

func Test_Expiration(t *testing.T) {
	s := New(Options{
		ExpiresIn: 3 * time.Second,
	})

	cookie := s.Set("abc", "mno")
	require.Equal(t, "mno", s.Get(cookie.Value))
	time.Sleep(3 * time.Second)
	require.Nil(t, s.Get(cookie.Value))
}
