package jim

import (
	"log"
	"net/http"
)

//
//  router
//  @Description: 保存路由和对应的处理方法
//
type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

//
//  addRoute
//  @Description: 将路由和对应的处理方法添加到Engine
//  @receiver engine
//  @param method
//  @param pattern
//  @param handle
//
func (router *router) addRoute(method string, pattern string, handle HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handle
}

//
//  handle
//  @Description:  执行路由的方法
//  @receiver router
//  @param c
//
func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
