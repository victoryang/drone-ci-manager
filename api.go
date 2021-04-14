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

	info := Rolling.GetBasicInfo(project)
	if info==nil {
		err := errors.New("Project not found in Rolling")
		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	err := CreateProject(project, info.GitURL)
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
	projects, err := GetAllProjects()
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

	info := Rolling.GetBasicInfo(project)
	if info==nil {
		err := errors.New("Project not found in Rolling")
		c.AbortWithStatusJSON(http.StatusBadRequest, formatErr(err))
		return
	}

	err := DeleteProject(project, info.GitURL)
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
		Info 		*Project 		`json:"projectInfo"`
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

	p := b.Info
	p.Project = project
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