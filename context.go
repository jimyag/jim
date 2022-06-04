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

	// 响应的信息
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
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
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
