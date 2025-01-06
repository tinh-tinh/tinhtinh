package compress_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
)

type BigStruct struct {
	ID          int64
	Name        string
	Description string
	Timestamp   string
	Tags        []string
	Data        []byte
	Metadata    map[string]interface{}
	Coordinates struct {
		Latitude  float64
		Longitude float64
	}
	Attributes struct {
		IsActive    bool
		AccessLevel int
		Notes       []string
	}
	Children []struct {
		ChildID   int
		ChildName string
		ChildData []byte
	}
}

func fillValue() BigStruct {
	bigStruct := BigStruct{
		ID:          123456789,
		Name:        "ExampleStruct",
		Description: "This is a large struct for testing purposes.",
		Timestamp:   time.Now().Format("2006-01-01"),
		Tags:        []string{"example", "testing", "golang"},
		Data:        []byte("Random binary data for testing."),
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": []string{"nested", "data"},
		},
		Children: []struct {
			ChildID   int
			ChildName string
			ChildData []byte
		}{
			{ChildID: 1, ChildName: "Child1", ChildData: []byte("Child1 data")},
			{ChildID: 2, ChildName: "Child2", ChildData: []byte("Child2 data")},
		},
	}

	// Fill coordinates and attributes
	bigStruct.Coordinates.Latitude = 37.7749
	bigStruct.Coordinates.Longitude = -122.4194
	bigStruct.Attributes.IsActive = true
	bigStruct.Attributes.AccessLevel = 5
	bigStruct.Attributes.Notes = []string{"Note1", "Note2", "Note3"}

	return bigStruct
}

func Test_Gzip(t *testing.T) {
	bigstruct := fillValue()
	val, err := compress.Encode(bigstruct, compress.Gzip)
	require.Nil(t, err)

	valRaw, err := compress.Decode[BigStruct](val, compress.Gzip)
	require.Nil(t, err)
	require.Equal(t, bigstruct, valRaw)
}

func Test_Flate(t *testing.T) {
	bigstruct := fillValue()
	val, err := compress.Encode(bigstruct, compress.Flate)
	require.Nil(t, err)

	valRaw, err := compress.Decode[BigStruct](val, compress.Flate)
	require.Nil(t, err)
	require.Equal(t, bigstruct, valRaw)
}

func Test_Zlib(t *testing.T) {
	bigstruct := fillValue()
	val, err := compress.Encode(bigstruct, compress.Zlib)
	require.Nil(t, err)

	valRaw, err := compress.Decode[BigStruct](val, compress.Zlib)
	require.Nil(t, err)
	require.Equal(t, bigstruct, valRaw)
}

func Test_Error(t *testing.T) {
	_, err := compress.Encode(nil, compress.Gzip)
	require.NotNil(t, err)

	_, err = compress.Encode(nil, compress.Flate)
	require.NotNil(t, err)

	_, err = compress.Encode(nil, compress.Zlib)
	require.NotNil(t, err)

	val, err := compress.Encode("test", "abc")
	require.NotNil(t, err)
	require.Nil(t, val)

	valRaw, err := compress.Decode[BigStruct](val, "abc")
	require.NotNil(t, err)
	require.Nil(t, valRaw)

	valRaw, err = compress.Decode[BigStruct](val, compress.Gzip)
	require.NotNil(t, err)
	require.Nil(t, valRaw)

	valRaw, err = compress.Decode[BigStruct](val, compress.Flate)
	require.NotNil(t, err)
	require.Nil(t, valRaw)

	valRaw, err = compress.Decode[BigStruct](val, compress.Zlib)
	require.NotNil(t, err)
	require.Nil(t, valRaw)
}

func Test_IsValidAlg(t *testing.T) {
	require.True(t, compress.IsValidAlg(compress.Gzip))
	require.True(t, compress.IsValidAlg(compress.Flate))
	require.True(t, compress.IsValidAlg(compress.Zlib))
	require.False(t, compress.IsValidAlg("abc"))
}
