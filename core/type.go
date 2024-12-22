package core

// Serialization
type Encode func(v interface{}) ([]byte, error)

type Decode func(data []byte, v interface{}) error
