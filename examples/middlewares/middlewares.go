package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jimyag/jim"
)

func onlyForV2() jim.HandleFunc {
	return func(c *jim.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := jim.New()
	r.Use(jim.Logger()) // global midlleware
	// jimyag@jimyagMac repo % curl http://localhost:9999/
	// <h1>Hello Jim</h1>
	// 2022/06/16 22:55:50 [200] / in 7.887µs
	r.GET("/", func(c *jim.Context) {
		c.HTML(http.StatusOK, "", "<h1>Hello Jim</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		// jimyag@jimyagMac repo % curl http://localhost:9999/v2/hello/jimyag
		// {"message":"Internal Server Error"}
		// 2022/06/16 22:56:24 [500] /v2/hello/jimyag in 88.357µs for group v2
		// 2022/06/16 22:56:24 [500] /v2/hello/jimyag in 106.744µs
		v2.GET("/hello/:name", func(c *jim.Context) {
			// expect /hello/jimyag
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9999")
}
