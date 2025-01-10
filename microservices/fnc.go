package microservices

import "fmt"

type ReqFnc func(event string, data interface{}, headers ...Header) error

func DefaultSend(c ClientProxy) ReqFnc {
	return func(event string, data interface{}, headers ...Header) error {
		err := c.Emit(event, Message{
			Type:    RPC,
			Event:   event,
			Headers: AssignHeader(c.Headers(), headers...),
			Data:    data,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Send mesage: %v for event: %s\n", data, event)
		return nil
	}
}

func DefaultPublish(c ClientProxy) ReqFnc {
	return func(event string, data interface{}, headers ...Header) error {
		err := c.Emit(event, Message{
			Type:    PubSub,
			Event:   event,
			Headers: AssignHeader(c.Headers(), headers...),
			Data:    data,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Publish mesage: %v for event: %s\n", data, event)
		return nil
	}
}
