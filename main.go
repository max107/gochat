package main

import (
	"log"
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/mimicloud/easyconfig"
	"github.com/flosch/pongo2"
	"github.com/robvdl/pongo2gin"
	"net/http"
	"path"
)

const configPath = "./config.json"

var config = struct {
	Path string `json:"path"`
	BaseUrl string `json:"base_url"`
	Listen string `json:"listen"`
}{}

func init() {
	// gin.SetMode(gin.ReleaseMode)
	// Read config file
	easyconfig.Parse(configPath, &config)
	if config.Path == "" {
		config.Path = "./repos/"
		easyconfig.Save(configPath, &config)
	}
	if config.Listen == "" {
		config.Listen = ":9123"
		easyconfig.Save(configPath, &config)
	}
	if config.BaseUrl == "" {
		config.BaseUrl = "git@192.168.0.249:/volume1/storage/repos"
		easyconfig.Save(configPath, &config)
	}
}

func ListDirectory() (map[string]string) {
	dirs := make(map[string]string)
	files, err := ioutil.ReadDir(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		dirs[f.Name()] = path.Join(config.BaseUrl, f.Name())
	}
	return dirs
}

type CreateForm struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func runCommand(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	return string(out), err
}

func CreateRepository(name string) error {
	repoPath := path.Join(config.Path, name) + ".git"
	if _, err := runCommand(fmt.Sprintf("mkdir %s", repoPath)); err == nil {
		if _, initErr := runCommand(fmt.Sprintf("git init --bare --shared %s", repoPath)); initErr == nil {
			return nil
		} else {
			return initErr
		}
	} else {
		return err
	}
}

func CheckRepository(name string) bool {
	repoPath := path.Join(config.Path, name)
	_, err := os.Stat(repoPath)
	return os.IsNotExist(err)
}

func main() {
	r := gin.Default()
	r.HTMLRender = pongo2gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", pongo2.Context{
			"repositories": ListDirectory(),
		})
	})
	r.POST("/", func(c *gin.Context) {
		var form CreateForm
		if c.Bind(&form) == nil {
			if form.Name != "" {
				if CheckRepository(form.Name) == false {
					c.HTML(http.StatusOK, "index.html", pongo2.Context{
						"repositories": ListDirectory(),
						"message": "Repository already exists",
					})
					return
				}

                if err := CreateRepository(form.Name); err == nil {
					c.HTML(http.StatusOK, "index.html", pongo2.Context{
						"repositories": ListDirectory(),
						"message": "ok",
					})
					return
	            } else {
					c.HTML(http.StatusOK, "index.html", pongo2.Context{
						"repositories": ListDirectory(),
						"message": fmt.Sprintf("Failed to create repository: %s", err),
					})
					return
	            }
			}
		} else {
			c.HTML(http.StatusOK, "index.html", pongo2.Context{
				"repositories": ListDirectory(),
			})
			return
		}
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/list", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"repositories": ListDirectory(),
			})
		})
		v1.POST("/create", func(c *gin.Context) {
			var form CreateForm
			if c.Bind(&form) == nil {
				if form.Name != "" {
					if CheckRepository(form.Name) == false {
						c.JSON(http.StatusOK, gin.H{"status": false, "message": "Repository already exists"})
						return
					}

					if err := CreateRepository(form.Name); err == nil {
						c.JSON(http.StatusOK, gin.H{"status": true, "message": "ok"})
					} else {
						c.JSON(http.StatusOK, gin.H{"status": false, "message": fmt.Sprintf("Failed to create repository: %s", err)})
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"status": false, "message": "Missing repository name"})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{"status": false, "message": "Failed to bind data"})
			}
		})
	}
	r.Run(config.Listen)
}
