/*
 * @Author: liziwei01
 * @Date: 2021-04-29 15:14:03
 * @LastEditors: liziwei01
 * @LastEditTime: 2021-05-30 02:35:36
 * @Description: http model
 * @FilePath: /github.com/liziwei01/go-liziwei01-library/model/ghttp/base.go
 */
package ghttp

import (
	"net/http"

	errBase "github.com/liziwei01/go-liziwei01-library/model/error"
)

type Ghttp interface {
	Request() **http.Request
	Response() *http.ResponseWriter
	Get(str string) string
	Post(str string) string
	Write(data interface{}, errno int, err error)
}

type ghttp struct {
	request  **http.Request
	response *http.ResponseWriter
}

/**
 * @description: return a default http model
 * @param {**http.Request} request
 * @param {*http.ResponseWriter} response
 * @return {*}
 */
func Default(request **http.Request, response *http.ResponseWriter) Ghttp {
	(*response).Header().Set("content-type", "text/json")
	(*response).Header().Set("Access-Control-Allow-Origin", "*")
	return &ghttp{
		request:  request,
		response: response,
	}
}

/**
 * @description: Request
 * @param {*}
 * @return {*}
 */
func (g *ghttp) Request() **http.Request {
	return g.request
}

/**
 * @description: Response
 * @param {*}
 * @return {*}
 */
func (g *ghttp) Response() *http.ResponseWriter {
	return g.response
}

/**
 * @description: get variables in url
 * @param {string} str
 * @return {*}
 */
func (g *ghttp) Get(str string) string {
	return (**g.Request()).URL.Query().Get(str)
}

/**
 * @description: get postform variables
 * @param {string} str
 * @return {*}
 */
func (g *ghttp) Post(str string) string {
	(**g.Request()).ParseForm()
	for k, v := range (**g.Request()).PostForm {
		if k == str {
			return v[0]
		}
	}
	return ""
}

/**
 * @description: write into response
 * @param {interface{}} data
 * @param {int} errno
 * @param {error} err
 * @return {*}
 */
func (g *ghttp) Write(data interface{}, errno int, err error) {
	switch errno {
	case errBase.ErrorNoSuccess:
		(*g.Response()).WriteHeader(200)
	case errBase.ErrorNoFailure:
		(*g.Response()).WriteHeader(404)
	case errBase.ErrorNoClient:
		(*g.Response()).WriteHeader(400)
	case errBase.ErrorNoServer:
		(*g.Response()).WriteHeader(500)
	case errBase.ErrorNoSign:
		(*g.Response()).WriteHeader(200)
	}
	if err != nil {
		(*g.Response()).Write(errBase.Marshal(data, errno, err.Error()))
		return
	}
	(*g.Response()).Write(errBase.Marshal(data, errno, ""))
}

/**
 * @description: customize AccessControlAllowOrigin
 * @param {string} allow
 * @return {*}
 */
func (g *ghttp) SetAccessControlAllowOrigin(allow string) {
	(*g.response).Header().Set("Access-Control-Allow-Origin", allow)
}
