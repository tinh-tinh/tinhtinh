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

func EncodeMessage(client ClientProxy, message Message) ([]byte, error) {
	payload, err := client.Serializer(message)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func DecodeMessage(svc Service, data []byte) Message {
	var message Message
	err := svc.Deserializer(data, &message)
	if err != nil {
		panic(err)
	}
	return message
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
