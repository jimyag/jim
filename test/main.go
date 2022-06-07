package main

import (
	"net/http"

	"github.com/jimyag/jim"
)

func main() {
	r := jim.New()
	//curl http://localhost:9999
	//<h1>Hello Jim</h1>
	r.GET("/", func(c *jim.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Jim</h1>")
	})

	// curl  http://localhost:9999/hello
	//hello, you're at /hello
	r.GET("/hello", func(c *jim.Context) {
		c.String(http.StatusOK, "hello, you're at %s\n", c.Path)
	})

	// curl  http://localhost:9999/hello/jimyag
	//hello jimyag, you're at /hello/jimyag
	r.GET("/hello/:name", func(c *jim.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	// curl http://localhost:9999/assets/css/geektutu.css
	//{"filepath":"css/geektutu.css"}
	r.GET("/assets/*filepath", func(c *jim.Context) {
		c.JSON(http.StatusOK, jim.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
