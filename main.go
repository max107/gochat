package main

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mimicloud/easyconfig"
	"github.com/robvdl/pongo2gin"
	"net/http"

	chat "./chat"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/doug-martin/goqu.v3"
	_ "gopkg.in/doug-martin/goqu.v3/adapters/sqlite3"
)

const configPath = "./config.json"

var config = struct {
	Listen string `json:"listen"`
	DbPath string `json:"db_path"`
}{}

var db *goqu.Database
var server *chat.Server

func init() {
	server = chat.NewServer()

	// gin.SetMode(gin.ReleaseMode)
	// Read config file
	easyconfig.Parse(configPath, &config)
	if config.Listen == "" {
		config.Listen = ":9123"
		easyconfig.Save(configPath, &config)
	}

	if config.DbPath == "" {
		config.DbPath = "./database.sqlite3"
		easyconfig.Save(configPath, &config)
	}

	sqliteDb, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		panic(err.Error())
	}

	db = goqu.New("sqlite3", sqliteDb)
}

func main() {
	r := gin.Default()
	// static files have higher priority over dynamic routes
	r.Use(static.Serve("/static", static.LocalFile("./static", false)))
	r.HTMLRender = pongo2gin.Default()
	r.Use(DummyMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", pongo2.Context{})
	})
	r.GET("/ws/:name", func(c *gin.Context) {
		if name := c.Param("name"); name != "" {
			server.UpdateWsHandler(name, c.Writer, c.Request)
		} else {
			c.Abort()
		}
	})

	r.Run(config.Listen)
}
