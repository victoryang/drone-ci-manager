package main

import (
	"net/http"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

//go:generate go-bindata -prefix "web/" -o=web.go -pkg=main web/...
//go:generate go-bindata -prefix "template/" -o=template/template.go -pkg=template template/...

func main(){
	router := gin.Default()
	router.Use(Logger())

	v1 := router.Group("/api/v1")
	for idx,_ :=range DroneServers {
		id := strconv.Itoa(idx)
		v1.POST("/drone/"+ id + "/buildyaml", gin.WrapH(NewYamlPlugin(idx)))
		v1.POST("/drone/" + id + "/webhook", gin.WrapH(NewWebhookPlugin(idx)))
	}

	projects := router.Group("/projects")
	projects.GET("/", getProjectList)
	projects.POST("/:project", createProject)
	projects.GET("/:project", getProjectInfo)
	projects.DELETE("/:project", deleteProject)
	projects.POST("/:project/scripts", createScripts)

	//้ๆๆไปถ
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	webFS := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
	router.StaticFS("/web", &webFS)
	router.Run(":9900")
}
