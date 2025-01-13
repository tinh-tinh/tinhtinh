package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
)

type EventType string

const (
	RPC    EventType = "rpc"
	PubSub EventType = "pubsub"
)

type Message struct {
	Type    EventType         `json:"type"`
	Event   string            `json:"event"`
	Headers map[string]string `json:"headers"`
	Data    interface{}       `json:"data"`
	Bytes   []byte            `json:"bytes"`
}

func EncodeMessage(c ClientProxy, message Message) ([]byte, error) {
	if c.Config().CompressAlg != "" {
		encoder, err := compress.Encode(message.Data, c.Config().CompressAlg)
		if err != nil {
			return nil, err
		}
		message.Data = nil
		message.Bytes = encoder
	}
	payload, err := c.Config().Serializer(message)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func DecodeMessage(c Service, data []byte) Message {
	var msg Message
	err := c.Config().Deserializer(data, &msg)
	if err != nil {
		panic(err)
	}
	if c.Config().CompressAlg != "" {
		decoder, err := compress.Decode(msg.Bytes, c.Config().CompressAlg)
		if err != nil {
			panic(err)
		}
		msg.Bytes = decoder
	}
	return msg
}

func AssignHeader(original Header, toMerge ...Header) Header {
	cloned := common.CloneMap(original)
	if len(toMerge) > 0 {
		for _, v := range toMerge {
			common.MergeMaps(cloned, v)
		}
	}

	return cloned
}
