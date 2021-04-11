package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// formatErr 错误流程返回错误给前端使用
func formatErr(err error) gin.H {
	return gin.H{"ErrMessage": fmt.Sprint(err)}
}

func createProject(c *gin.Context) {
	project := c.Param("project")

	err := CreateWorkingDir(project)
	if err!=nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func getProjectList(c *gin.Context) {
	projects, err := GetProjectsFromDir()
	if err!=nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	c.JSON(http.StatusOK, projects)
}

func getProjectInfo(c *gin.Context) {
	project := c.Param("project")

	info := Rolling.GetBasicInfo(project)
	if info==nil {
		err := errors.New("Project not found in Rolling")
		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	c.JSON(http.StatusOK, info)
}

func deleteProject(c *gin.Context) {
	project := c.Param("project")

	err := DeleteProjectFromDir(project)
	if err!=nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func createScripts(c *gin.Context) {
	project := c.Param("project")

	type Body struct {
		Info 		*ProjectInfo 	`json:"projectInfo"`
		Envs 		[]string 		`json:"envs"`
	}

	b := &Body{}
	err := c.BindJSON(&b)
	if err!=nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	p := &Project{
		Project: project,
		Info: b.Info,
	}
	err = p.generateFiles(b.Envs)
	if err!=nil {
		defer func() {
			c.Error(err)
		}()

		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}