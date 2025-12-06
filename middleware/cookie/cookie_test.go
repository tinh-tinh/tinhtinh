package cookie_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/cookie"
)

func Test_SecureCookie(t *testing.T) {
	assert.Equal(t, "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92", cookie.Encode([]byte("oievrervierjvoiejrojojov")))
	decode, err := cookie.Decode("b2lldnJlcnZpZXJqdm9pZWpyb2pvam92")
	assert.Nil(t, err)
	assert.Equal(t, "oievrervierjvoiejrojojov", string(decode))

	_, err = cookie.Decode("Tôi Tích Ta Tu Tiên")
	assert.NotNil(t, err)

	sCookie := &cookie.SecureCookie{Key: "add"}
	_, err = sCookie.Encrypt("abc")
	assert.NotNil(t, err)

	_, err = sCookie.Decrypt("avv")
	assert.NotNil(t, err)

	sCookie2 := &cookie.SecureCookie{Key: "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92"}
	_, err = sCookie2.Decrypt("Tôi Tích Ta Tu Tiên")
	assert.NotNil(t, err)

	sCookie3 := &cookie.SecureCookie{Key: "b2lldnJlcnZpZXJqdm9pZWpyb2pvam92"}
	encrypted, err := sCookie3.Encrypt("Hello World!")
	assert.Nil(t, err)

	decrypted, err := sCookie3.Decrypt(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", decrypted)

	// Test with another key
	sCookie4 := &cookie.SecureCookie{Key: "c2VjdXJla2V5MTIzNDU2Nzg5MA==rttt"}
	encrypted2, err := sCookie4.Encrypt("Xin chào thế giới!")
	assert.Nil(t, err)

	decrypted2, err := sCookie4.Decrypt(encrypted2)
	assert.Nil(t, err)
	assert.Equal(t, "Xin chào thế giới!", decrypted2)

	// Ensure that decrypting with a different key fails
	_, err = sCookie3.Decrypt(encrypted2)
	assert.NotNil(t, err)
}

func Test_Pkcs7PadUnpad(t *testing.T) {
	data := []byte("Yêu thích tiên hiệp")
	blockSize := 16

	padded := cookie.Pkcs7Pad(data, blockSize)
	unpadded, err := cookie.Pkcs7Unpad(padded, blockSize)
	assert.Nil(t, err)
	assert.Equal(t, data, unpadded)

	// Test unpad with invalid padding
	invalidPadded := append(data, 0x05, 0x05, 0x05)
	_, err = cookie.Pkcs7Unpad(invalidPadded, blockSize)
	assert.NotNil(t, err)

	// Test unpad with empty data
	_, err = cookie.Pkcs7Unpad([]byte{}, blockSize)
	assert.NotNil(t, err)

	// Test unpad with invalid padding
	invalidPadded2 := append(padded, 0)
	_, err = cookie.Pkcs7Unpad(invalidPadded2, (blockSize*2)+1)
	assert.NotNil(t, err)

	invalidPadded3 := append(padded, 0x5)
	_, err = cookie.Pkcs7Unpad(invalidPadded3, (blockSize*2)+1)
	assert.NotNil(t, err)
}
