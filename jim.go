package jim

import (
	"log"
	"net/http"
)

// HandleFunc jim 对请求进行处理的方法
type HandleFunc func(ctx *Context)

//
// Engine 实现了 ServeHTTP 的接口
type Engine struct {
	*RouterGroup
	// 保存路由和对应的处理方法
	router *router
	groups []*RouterGroup // 存储所有的 groups
}

//
// New
//  @Description:  创建一个jim.Engine
//  @return *Engine
//
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
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

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc // 支持中间件
	parent      *RouterGroup // 支持嵌套
	engine      *Engine      // 所有的 Groups 共享一个 Engine 实例
}

//
// Group
//  @Description: Group 定义为创建一个新的 RouterGroup
// 记住所有 Group 共享相同的 Engine 实例
//  @receiver group
//  @param prefix
//  @return *RouterGroup
//
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//
// GET
//  @Description:  GET defines the method to add GET request
//  @receiver group
//  @param pattern
//  @param handler
//
func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

//
// POST
//  @Description: POST defines the method to add POST request
//  @receiver group
//  @param pattern
//  @param handler
//
func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}
