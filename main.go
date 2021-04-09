package main

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

//go:generate go-bindata -prefix "web/" -o=web.go -pkg=main web/...

func main(){
	router := gin.Default()
	router.Use(Logger())

	v1 := router.Group("/api/v1")
	v1.POST("/projects/:project", createRollingProject)
	v1.POST("/drone/build-yaml", gin.WrapH(NewYamlPlugin()))
	v1.POST("/drone/webhook", gin.WrapH(NewWebhookPlugin()))
	//v1.POST("/drone/registry-info", gin.WrapH(NewRegistryPlugin()))

	projects := router.Group("/projects")
	projects.GET("/", getProjectList)
	projects.GET("/:project", getProjectInfo)
	projects.POST("/:project/scripts", createScripts)
	projects.POST("/:project", createProject)
	projects.DELETE("/:project", deleteProject)

	//静态文件
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	webFS := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
	router.StaticFS("/web", &webFS)
	router.Run(":5000")
}