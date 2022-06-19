package main

import (
	"net/http"

	"github.com/jimyag/jim"
)

func main() {
	r := jim.New()
	r.GET("/index", func(c *jim.Context) {
		c.HTML(http.StatusOK, "", "<h1>Index Page</h1>")
	})

	// curl http://localhost:9999/v1/hello?name=jimyag
	//hello jimyag, you're at /v1/hello
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *jim.Context) {
			c.HTML(http.StatusOK, " ", "<h1>Hello Jim</h1>")
		})

		v1.GET("/hello", func(c *jim.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	// curl "http://localhost:9999/v2/hello/jimyag"
	//hello jimyag, you're at /v2/hello/jimyag
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *jim.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *jim.Context) {
			c.JSON(http.StatusOK, jim.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
