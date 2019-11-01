package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rakyll/statik/fs"
	_ "github.com/theoremoon/kaidsuka/statikspa/statik"
)

//go:generate statik -src ./dist -f

func main() {
	fs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/*", func(c echo.Context) error {
		res := c.Response()
		writer := res.Writer

		recorder := httptest.NewRecorder()
		res.Writer = recorder

		handler := echo.WrapHandler(http.FileServer(fs))
		err := handler(c)
		if err != nil {
			return err
		}
		res.Writer = writer

		if res.Status == http.StatusNotFound {
			f, err := fs.Open("/index.html")
			if err != nil {
				log.Println(err)
				return nil
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			return c.HTMLBlob(http.StatusOK, content)

		} else {
			writer.Write(recorder.Body.Bytes())
		}

		return nil
	})
	e.GET("/api/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Key": "Value",
		})
	})

	e.Logger.Fatal(e.Start(":8888"))
}
