package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
)

type Alg string

const (
	Gzip  Alg = "gzip"
	Flate Alg = "flate"
	Zlib  Alg = "zlib"
)

func Encode(data interface{}, alg Alg, levels ...int) ([]byte, error) {
	valBytes, err := ToBytes(data)
	if err != nil {
		return nil, err
	}

	if len(levels) == 0 {
		levels = []int{gzip.DefaultCompression}
	}

	var buf bytes.Buffer
	switch alg {
	case Gzip:
		gzipWriter, err := gzip.NewWriterLevel(&buf, levels[0])
		if err != nil {
			return nil, err
		}
		_, err = gzipWriter.Write(valBytes)
		if err != nil {
			return nil, err
		}
		gzipWriter.Close()
	case Flate:
		flateWriter, err := flate.NewWriter(&buf, levels[0])
		if err != nil {
			return nil, err
		}
		_, err = flateWriter.Write(valBytes)
		if err != nil {
			return nil, err
		}
		flateWriter.Close()
	case Zlib:
		zlibWriter, err := zlib.NewWriterLevel(&buf, levels[0])
		if err != nil {
			return nil, err
		}
		_, err = zlibWriter.Write(valBytes)
		if err != nil {
			return nil, err
		}
		zlibWriter.Close()
	default:
		return nil, errors.New("unknown compression algorithm")
	}

	return buf.Bytes(), nil
}

func Decode(data []byte, alg Alg) ([]byte, error) {
	buf := bytes.NewReader(data)
	var decompressed bytes.Buffer

	switch alg {
	case Gzip:
		gzipReader, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()

		_, err = io.Copy(&decompressed, gzipReader)
		if err != nil {
			return nil, err
		}
	case Flate:
		flateReader := flate.NewReader(buf)
		defer flateReader.Close()

		_, err := io.Copy(&decompressed, flateReader)
		if err != nil {
			return nil, err
		}
	case Zlib:
		zlibReader, err := zlib.NewReader(buf)
		if err != nil {
			return nil, err
		}
		defer zlibReader.Close()

		_, err = io.Copy(&decompressed, zlibReader)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown compression algorithm")
	}

	return decompressed.Bytes(), nil
}

func DecodeMarshall[T any](data []byte, alg Alg) (T, error) {
	decompress, err := Decode(data, alg)
	if err != nil {
		return *new(T), err
	}
	return FromBytes[T](decompress)
}
