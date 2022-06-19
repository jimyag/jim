package jim

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

//
// Context
//  @Description: 定义上下文，可以减少大量重复的代码
//
type Context struct {
	// 原始的对象
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求相关的信息
	Path   string
	Method string
	Params map[string]string // 解析的参数 c.Param("lang")的方式获取到对应的值。

	// 响应的信息
	StatusCode int

	// 中间件
	handlers []HandleFunc
	index    int // 记录当前到了第几个中间件

	engine *Engine
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

// Next 继续执行下一个中间件 当调用 Next 的时候，控制权交给了下一个中间件，
// 		直到调用了最后一个中间件，然后从后往前
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//
// PostForm
//  @Description: 访问 PostForm 参数
//  @receiver c
//  @param key
//  @return string
//
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//
// Query
//  @Description: 访问 Query 参数
//  @receiver c
//  @param key
//  @return string
//
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//
// Status
//  @Description: 设置 Status
//  @receiver c
//  @param code
//
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//
// SetHeader
//  @Description: 设置 header
//  @receiver c
//  @param key
//  @param value
//
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//
//  String
//  @Description: 构造 String 响应的方法
//  @receiver c
//  @param code
//  @param format
//  @param values
//
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//
// JSON
//  @Description: 构造 JSON 响应的方法
//  @receiver c
//  @param code
//  @param obj
//
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//
// Data
//  @Description: 构造 Data 响应方法
//  @receiver c
//  @param code
//  @param data
//
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//
// HTML
//  @Description:  构造 HTML 响应方法
//  @receiver c
//  @param code
//  @param html
//
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

//
// Param
//  @Description: 获得路由中的参数
//	路由是 /p/:lang/s 请求的路由是 /p/aa/s
//	此时 c.Params {lang:"aa"} 调用 c.Param("lang") 可以得到 "aa"
//  @receiver c
//  @param key
//  @return string
//
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Fail
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
