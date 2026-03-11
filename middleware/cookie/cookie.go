package cookie

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

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
	key := []byte(s.Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	// CBC needs a block-sized IV; use a new random IV per encryption.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("iv: %w", err)
	}

	plain := []byte(text)
	padded := Pkcs7Pad(plain, block.BlockSize())

	cipherText := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText, padded)

	// Prefix IV so it can be used for decryption.
	out := append(iv, cipherText...)

	return Encode(out), nil
}

func (s *SecureCookie) Decrypt(text string) (string, error) {
	key := []byte(s.Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	cipherTextWithIV, err := Decode(text)
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	bs := block.BlockSize()
	if len(cipherTextWithIV) < bs || len(cipherTextWithIV)%bs != 0 {
		return "", fmt.Errorf("ciphertext too short or misaligned")
	}

	iv := cipherTextWithIV[:bs]
	cipherText := cipherTextWithIV[bs:]

	plainPadded := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainPadded, cipherText)

	plain, err := Pkcs7Unpad(plainPadded, bs)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// Pkcs7Pad pads data to a multiple of blockSize.
func Pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - (len(data) % blockSize)
	return append(data, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

func Pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}
	for _, v := range data[len(data)-pad:] {
		if int(v) != pad {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}
	return data[:len(data)-pad], nil
}
