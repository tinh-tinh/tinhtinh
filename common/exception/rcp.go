package exception

import "encoding/json"

type Rpc struct {
	Code string
	Msg  string
}

func ThrowRpc(msg string, code ...string) Rpc {
	rpc := Rpc{Msg: msg}
	if len(code) > 0 {
		rpc.Code = code[0]
	}
	return rpc
}

func (e Rpc) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func AdapterHttpRpc(err error) Rpc {
	var e Rpc
	er := json.Unmarshal([]byte(err.Error()), &e)
	if er != nil {
		return Rpc{Msg: err.Error()}
	}
	return e
}
