package error

import (
	"encoding/json"

	errLib "github.com/liziwei01/go-liziwei01-library/library/error"
)

const (
	ErrorMsgSuccess = "success"
	ErrorMsgFailure = "failure"
	ErrorMsgClient  = "params check failed"
	ErrorMsgServer  = "server failed"
	ErrorMsgSign    = "sign check failed"
	ErrorNoSuccess  = 0
	ErrorNoFailure  = -1
	ErrorNoClient   = -2
	ErrorNoServer   = -3
	ErrorNoSign     = -4
)

func Marshal(data interface{}, errno int, errmsg string) []byte {
	switch errno {
	case ErrorNoSuccess:
		errmsg = ErrorMsgSuccess
	case ErrorNoFailure:
		errmsg = ErrorMsgFailure
	case ErrorNoClient:
		errmsg = ErrorMsgClient
	case ErrorNoServer:
		errmsg = ErrorMsgServer
	case ErrorNoSign:
		errmsg = ErrorMsgSign
	}
	e := errLib.New(errLib.ErrNo(errno), errLib.ErrMsg(errmsg))
	errLib.Append(e, errLib.ErrMsg(errmsg))
	ret, _ := json.Marshal(map[string]interface{}{
		"data":   data,
		"errmsg": e.ErrMsg(),
		"errno":  e.ErrNo(),
	})
	return ret
}

func MarshalData(data interface{}) []byte {
	e := errLib.NewDefault()
	ret, _ := json.Marshal(map[string]interface{}{
		"data":   data,
		"errmsg": e.ErrMsg(),
		"errno":  e.ErrNo(),
	})
	return ret
}
