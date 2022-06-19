package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jimyag/jim"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := jim.New()
	r.Use(jim.Logger())

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "jimyag", Age: 20}
	stu2 := &student{Name: "jack", Age: 22}
	r.GET("/", func(c *jim.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *jim.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", jim.H{
			"title":  "Jim",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *jim.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", jim.H{
			"title": "jimyag",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
