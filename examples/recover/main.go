package main

import (
	"net/http"

	"github.com/jimyag/jim"
)

func main() {
	r := jim.Default()
	r.GET("/", func(c *jim.Context) {
		c.String(http.StatusOK, "Hello Jimyag\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *jim.Context) {
		names := []string{"jimyag"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
