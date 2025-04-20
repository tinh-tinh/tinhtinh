package compress_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
)

func Test_Byte(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := &Person{
		Name: "John",
		Age:  30,
	}

	data, err := compress.ToBytes(person)
	require.Nil(t, err)

	val, err := compress.FromBytes[Person](data)
	require.Nil(t, err)

	require.Equal(t, person.Age, val.Age)
	require.Equal(t, person.Name, val.Name)

	_, err = compress.ToBytes(nil)
	require.NotNil(t, err)

	_, err = compress.FromBytes[Person]([]byte("test"))
	require.NotNil(t, err)
}
