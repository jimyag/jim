package jim

import (
	"net/http"
)

// HandleFunc jim 对请求进行处理的方法
type HandleFunc func(ctx *Context)

//
// Engine 实现了 ServeHTTP 的接口
type Engine struct {
	// 保存路由和对应的处理方法
	router *router
}

//
// New
//  @Description:  创建一个jim.Engine
//  @return *Engine
//
func New() *Engine {
	return &Engine{router: newRouter()}
}

//
//  addRoute
//  @Description: 将路由和对应的处理方法添加到Engine
//  @receiver engine
//  @param method
//  @param pattern
//  @param handle
//
func (engine *Engine) addRoute(method string, pattern string, handle HandleFunc) {
	engine.router.addRoute(method, pattern, handle)
}

//
// GET
//  @Description: GET 请求
//  @receiver engine
//  @param pattern
//  @param handleFunc
//
func (engine *Engine) GET(pattern string, handleFunc HandleFunc) {
	engine.addRoute("GET", pattern, handleFunc)
}

//
// POST
//  @Description: POST 请求
//  @receiver engine
//  @param pattern
//  @param handleFunc
//
func (engine *Engine) POST(pattern string, handleFunc HandleFunc) {
	engine.addRoute("POST", pattern, handleFunc)
}

//
// Run
//  @Description: 启动 HTTP server
//  @receiver engine
//  @param addr
//  @return error
//
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

//
//  ServeHTTP
//  @Description:
//  @receiver engine
//  @param w
//  @param r
//
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}
