package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/middleware/session"
)

func Test_Session(t *testing.T) {
	s := session.New(session.Options{
		GeneratorID: func() string {
			return time.Now().String()
		},
	})

	cookie := s.Set("abc", "mno")
	require.Equal(t, "mno", s.Get(cookie.Value))
}

func Test_Expiration(t *testing.T) {
	s := session.New(session.Options{
		ExpiresIn: 1 * time.Second,
	})

	cookie := s.Set("abc", "mno")
	require.Equal(t, "mno", s.Get(cookie.Value))
	time.Sleep(3 * time.Second)
	require.Nil(t, s.Get(cookie.Value))
}
