/*
 * @Author: 		liziwei01
 * @Date: 			2021-04-19 15:00:00
 * @LastEditTime: 	2021-04-19 15:00:00
 * @LastEditors: 	liziwei01
 * @Description: 	start http server and start listening
 * @FilePath: 		github.com/liziwei01/go-liziwei01-library/httpapi/httpapi.go
 */
package httpapi

import (
	"io"
	"net/http"
)

/**
 * @description: start http server and start listening
 * @param {*}
 * @return {*}
 */
func InitRouters() {
	// init routers
	// Routers.Init()

	// safe router
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Hello! THis is Ziwei. Welcome to my website!")
	})
}
