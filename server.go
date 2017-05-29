package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

func GetFile(c echo.Context) (err error) {
	filename := c.Param("filename")
	if filename == "" {
		return c.String(http.StatusBadRequest, "filename not input.")
	}
	log.Infof("Get file request %s", filename)
	f, err := os.Open("/home/frank/Downloads/test/" + filename + ".pdf")
	if err != nil {
		log.Errorf("Get file error. %s", err.Error())
		return err
	}
	return c.Stream(http.StatusOK, "Content-type:application/pdf", f)
}

func main() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	e.Static("/", "public")
	e.File("/viewer", "public/web/viewer.html")
	e.File("/", "public/index.html")
	e.GET("/file/:filename", GetFile)
	e.Logger.Fatal(e.Start(":1323"))
}
