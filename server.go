package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/eternnoir/webpdfviewer/db"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	FlagDataFolder = ""
	FlagServerPort = ""
	WelcomePdf     = "HELLO"
)

var dbmgr *db.DbMgr

type FileHistory struct {
	Filename string
	Record   []db.ViewRecord
}

func GetFile(c echo.Context) (err error) {
	filename := c.Param("filename")
	if filename == "" {
		return c.String(http.StatusBadRequest, "filename not input.")
	}
	log.Infof("Get file request %s", filename)
	fullfileName, err := tryFindFile(filename)
	if err != nil {
		log.Errorf("Filename %s not found. %s", filename, err.Error())
	}

	f, err := os.Open(fullfileName)
	if err != nil {
		log.Errorf("Get file error. %s", err.Error())
		return err
	}
	if filename != WelcomePdf {
		go dbmgr.InsertRecord(filename)
	}
	return c.Stream(http.StatusOK, "Content-type:application/pdf", f)
}

func GetFileHistory(c echo.Context) (err error) {
	filename := c.QueryParam("filename")
	if filename == "" {
		return c.Render(http.StatusOK, "query", "World")
	}
	return c.Render(http.StatusOK, "query", "World")
}

func tryFindFile(filename string) (string, error) {
	p1 := FlagDataFolder + filename + ".pdf"
	p2 := FlagDataFolder + filename + ".PDF"
	if _, err := os.Stat(p1); err == nil {
		return p1, nil
	}
	if _, err := os.Stat(p2); err == nil {
		return p2, nil
	}
	return "", fmt.Errorf("%s file %s.pdf or .PDF not found", filename, filename)
}

func start(c *cli.Context) error {
	if FlagDataFolder == "" {
		panic("Data folder must set.")
	}
	log.Infof("Use datafolder %s", FlagDataFolder)
	dm, err := db.NewDbMgr("./w.db")
	if err != nil {
		return err
	}
	dbmgr = dm

	t := &Template{
		templates: template.Must(template.ParseGlob("public/template/*.html")),
	}

	log.Infof("Ts %#v", t)

	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	e.Static("/", "public")
	e.File("/viewer", "public/web/viewer.html")
	e.File("/", "public/index.html")
	e.Renderer = t
	e.GET("/file/:filename", GetFile)
	e.GET("/query", GetFileHistory)
	e.Logger.Fatal(e.Start(FlagServerPort))
	return nil
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "folder, f",
			Value:       "./data",
			Usage:       "data path",
			Destination: &FlagDataFolder,
		},
		cli.StringFlag{
			Name:        "addr, a",
			Value:       ":8080",
			Usage:       "server listen address. Ex :8080",
			Destination: &FlagServerPort,
		},
	}

	app.Action = start
	app.Run(os.Args)
}
