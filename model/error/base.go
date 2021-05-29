/*
 * @Author: liziwei01
 * @Date: 2021-04-29 15:14:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2021-05-30 02:44:12
 * @Description: error model
 * @FilePath: /github.com/liziwei01/go-liziwei01-library/model/error/base.go
 */
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

/**
* @description: could use error number to auto fill message, but also customize
* 				return a json like
* 				{
* 				"data":   data,
* 				"errmsg": e.ErrMsg(),
* 				"errno":  e.ErrNo(),
* 				}
* @param {interface{}} data
* @param {int} errno
* @param {string} errmsg
* @return {*}
 */
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
	default:
		// do nothing
	}
	e := errLib.New(errLib.ErrNo(errno), errLib.ErrMsg(errmsg))
	ret, _ := json.Marshal(map[string]interface{}{
		"data":   data,
		"errmsg": e.ErrMsg(),
		"errno":  e.ErrNo(),
	})
	return ret
}

/**
* @description: if no error, use this to write data to response
* 				return a json like
* 				{
* 				"data":   data,
* 				"errmsg": "success",
* 				"errno":  0,
* 				}
 * @param {interface{}} data
 * @return {*}
*/
func MarshalData(data interface{}) []byte {
	return Marshal(data, 0, "")
}
