package ghttp

import (
	"net/http"

	errBase "github.com/liziwei01/go-liziwei01-library/model/error"
)

type Ghttp interface {
	Request() **http.Request
	Response() *http.ResponseWriter
}

type ghttp struct {
	request  **http.Request
	response *http.ResponseWriter
}

func Default(request **http.Request, response *http.ResponseWriter) Ghttp {
	(*response).Header().Set("content-type", "text/json")
	(*response).Header().Set("Access-Control-Allow-Origin", "*")
	return &ghttp{
		request:  request,
		response: response,
	}
}

func (g *ghttp) Request() **http.Request {
	return g.request
}

func (g *ghttp) Response() *http.ResponseWriter {
	return g.response
}

func Get(g Ghttp, str string) string {
	return (**g.Request()).URL.Query().Get(str)
}

func Post(g Ghttp, str string) string {
	(**g.Request()).ParseForm()
	for k, v := range (**g.Request()).PostForm {
		if k == str {
			return v[0]
		}
	}
	return ""
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