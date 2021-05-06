package ghttp

import (
	"net/http"

	errBase "github.com/go-liziwei01-library/modules/erg3020/model/error"
)

type (
	Request  *http.Request
	Response http.ResponseWriter
)

type Ghttp interface {
	Request() *Request
	Response() *Response
}

type ghttp struct {
	request  *Request
	response *Response
}

func Default(request *Request, response *Response) Ghttp {
	(*response).Header().Set("content-type", "text/json")
	(*response).Header().Set("Access-Control-Allow-Origin", "*")
	return &ghttp{
		request:  request,
		response: response,
	}
}

func (g *ghttp) Request() *Request {
	return g.request
}

func (g *ghttp) Response() *Response {
	return g.response
}

func Write(g Ghttp, data interface{}, errno int, err error) {
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
