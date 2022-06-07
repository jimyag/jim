package jim

import (
	"log"
	"net/http"
	"strings"
)

//
//  router
//  @Description: 保存路由和对应的处理方法
//
type router struct {
	roots    map[string]*node      //存储每种请求方式的Trie 树根节点
	handlers map[string]HandleFunc // 保存路由的处理方法
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandleFunc),
		roots:    make(map[string]*node),
	}
}

//
//  parsePattern
//  @Description: 只会匹配一个 *
//	/p/a/b/c  [p,a,b,c]
//	/p/a/*/c  [p,a,*]
//	/p/a/*/d/* [p,a,*]
//  @param pattern
//  @return []string
//
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
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
	// 划分路由
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	// 根据请求方法类型进行区分
	_, ok := router.roots[method]
	if !ok {
		router.roots[method] = &node{}
	}

	//  把 pattern 插入
	router.roots[method].insert(pattern, parts, 0)
	router.handlers[key] = handle
}

//
//  getRouter
//  @Description: 解析 : 和 * 两种匹配符的参数，返回一个 map 。
// 	例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}
//	/static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}。
//  @receiver router
//  @param method
//  @param path
//  @return *node
//  @return map[string]string
//
func (router *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//
//  handle
//  @Description:  执行路由的方法
//  @receiver router
//  @param c
//
func (router *router) handle(c *Context) {
	n, params := router.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		router.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
