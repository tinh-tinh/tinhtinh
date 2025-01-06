package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/gob"
	"errors"
	"io"
)

type Alg string

const (
	Gzip  Alg = "gzip"
	Flate Alg = "flate"
	Zlib  Alg = "zlib"
)

func IsValidAlg(alg Alg) bool {
	return alg == Gzip || alg == Flate || alg == Zlib
}

func Encode(plain interface{}, alg Alg, level ...int) ([]byte, error) {
	var trans bytes.Buffer
	enc := gob.NewEncoder(&trans)
	if err := enc.Encode(plain); err != nil {
		return nil, err
	}
	data := trans.Bytes()

	var buf bytes.Buffer
	if len(level) == 0 {
		level = []int{-1}
	}
	switch alg {
	case Gzip:
		writer, err := gzip.NewWriterLevel(&buf, level[0])
		if err != nil {
			return nil, err
		}
		_, err = writer.Write(data)
		if err != nil {
			return nil, err
		}
		writer.Close()
		return buf.Bytes(), nil
	case Flate:
		writer, err := flate.NewWriter(&buf, level[0])
		if err != nil {
			return nil, err
		}
		_, err = writer.Write(data)
		if err != nil {
			return nil, err
		}
		writer.Close()
		return buf.Bytes(), nil
	case Zlib:
		writer, err := zlib.NewWriterLevel(&buf, level[0])
		if err != nil {
			return nil, err
		}
		_, err = writer.Write(data)
		if err != nil {
			return nil, err
		}
		writer.Close()
		return buf.Bytes(), nil
	default:
		return nil, errors.New("invalid compress algorithm")
	}
}

func Decode[M any](data []byte, alg Alg) (interface{}, error) {
	buf := bytes.NewReader(data)
	var decompressed bytes.Buffer

	switch alg {
	case Gzip:
		reader, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		if _, err := io.Copy(&decompressed, reader); err != nil {
			return nil, err
		}
	case Flate:
		reader := flate.NewReader(buf)
		defer reader.Close()
		if _, err := io.Copy(&decompressed, reader); err != nil {
			return nil, err
		}
	case Zlib:
		reader, err := zlib.NewReader(buf)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		if _, err := io.Copy(&decompressed, reader); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid compress algorithm")
	}

	trans := bytes.NewBuffer(decompressed.Bytes())
	dec := gob.NewDecoder(trans)
	var result M
	for {
		err := dec.Decode(&result)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
