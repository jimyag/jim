package main

import (
	"net/http"

	"github.com/jimyag/jim"
)

func main() {
	r := jim.New()

	// curl -i http://localhost:9999/
	//HTTP/1.1 200 OK
	//Content-Type: text/html
	//Date: Sat, 04 Jun 2022 17:09:55 GMT
	//Content-Length: 18
	//
	//<h1>Hello Jim</h1>
	r.GET("/", func(c *jim.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Jim</h1>")
	})

	//curl "http://localhost:9999/hello?name=jimyag"
	//hello jimyag, you're at /hello
	r.GET("/hello", func(c *jim.Context) {
		// expect /hello?name=jimyag
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	//curl "http://localhost:9999/login" -X POST -d 'username=jimyag&password=jimyag'
	//{"password":"jimyag","username":"jimyag"}
	r.POST("/login", func(c *jim.Context) {
		c.JSON(http.StatusOK, jim.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
