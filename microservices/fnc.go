package microservices

type ReqFnc func(event string, data interface{}, headers ...Header) error

func DefaultSend(c ClientProxy) ReqFnc {
	return func(event string, data interface{}, headers ...Header) error {
		err := c.Emit(event, Message{
			Type:    RPC,
			Event:   event,
			Headers: AssignHeader(c.Config().Header, headers...),
			Data:    data,
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func DefaultPublish(c ClientProxy) ReqFnc {
	return func(event string, data interface{}, headers ...Header) error {
		err := c.Emit(event, Message{
			Type:    PubSub,
			Event:   event,
			Headers: AssignHeader(c.Config().Header, headers...),
			Data:    data,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
