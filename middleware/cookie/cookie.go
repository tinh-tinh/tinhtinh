package cookie

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}

type Options struct {
	Key string
}

type Key string

const SIGNED_COOKIE Key = "SignedCookie"

func Handler(opt Options) func(http.Handler) http.Handler {
	s := &SecureCookie{Key: opt.Key}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), SIGNED_COOKIE, s))
			next.ServeHTTP(w, r)
		})
	}
}

type SecureCookie struct {
	Key string
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (s *SecureCookie) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(s.Key))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func (s *SecureCookie) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(s.Key))
	if err != nil {
		return "", err
	}

	cipherText, err := Decode(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
