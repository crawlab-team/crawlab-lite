package controllers

import (
	"crawlab-lite/forms"
	"crawlab-lite/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

func QueryProjectList(c *gin.Context) {
	var page forms.PageForm

	if err := c.ShouldBind(&page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if total, projects, err := services.GetProjectList(page.PageNum, page.PageSize); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccessList(c, total, projects)
	}
}

func QueryProject(c *gin.Context) {
	name := c.Param("name")

	if project, err := services.GetProjectByName(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		if project == nil {
			HandleError(http.StatusNotFound, c, errors.New("project not found"))
			return
		}
		HandleSuccess(c, project)
	}
}

func CreateProject(c *gin.Context) {
	var form forms.ProjectForm

	if err := c.ShouldBind(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 正则校验项目名称
	if ok, err := regexp.MatchString("[\\w_-]", form.Name); err != nil || ok == false {
		HandleError(http.StatusBadRequest, c, errors.New("invalid project name"))
	}

	if res, err := services.SaveProject(form); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func DeleteProject(c *gin.Context) {
	name := c.Param("name")

	if res, err := services.RemoveProject(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func QueryProjectVersionList(c *gin.Context) {
	name := c.Param("name")

	if res, err := services.GetProjectVersionList(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func UploadProjectVersion(c *gin.Context) {
	var form forms.ProjectUploadForm

	if err := c.ShouldBind(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 如果不为 zip 文件，返回错误
	if !strings.HasSuffix(form.File.Filename, ".zip") {
		HandleError(http.StatusBadRequest, c, errors.New("invalid zip file"))
		return
	}

	projectName := c.Param("name")
	if res, err := services.SaveProjectVersion(projectName, form); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func DeleteProjectVersion(c *gin.Context) {
	name := c.Param("name")
	versionId := c.Param("versionId")

	if res, err := services.RemoveProjectVersion(name, versionId); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}
