/*
 * @Author: 		liziwei01
 * @Date: 			2021-04-19 15:00:00
 * @LastEditTime: 	2021-04-19 15:00:00
 * @LastEditors: 	liziwei01
 * @Description: 	启动http服务器并开始监听
 * @FilePath: 		github.com/Bill-xyz/go-liziwei01-library/httpapi/httpapi.go
 */
package httpapi

import (
	"io"
	"net/http"
)

/**
 * @description: 后台启动路由分发
 * @param {*}
 * @return {*}
 */
func InitRouters() {
	// init routers
	// Routers.Init()

	// 兜底路由
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Hello! THis is Ziwei. Welcome to my website!")
	})
}
