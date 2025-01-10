package microservices

import "github.com/tinh-tinh/tinhtinh/v2/common"

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
}

func EncodeMessage(c ClientProxy, message Message) ([]byte, error) {
	payload, err := c.Serializer(message)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func DecodeMessage(c Service, data []byte) Message {
	var msg Message
	err := c.Deserializer(data, &msg)
	if err != nil {
		panic(err)
	}
	return msg
}

func AssignHeader(orignal Header, toMerge ...Header) Header {
	cloned := common.CloneMap(orignal)
	if len(toMerge) > 0 {
		for _, v := range toMerge {
			common.MergeMaps(cloned, v)
		}
	}

	return cloned
}
