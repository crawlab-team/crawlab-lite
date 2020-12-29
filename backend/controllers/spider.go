package controllers

import (
	"crawlab-lite/forms"
	"crawlab-lite/services"
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
)

func GetSpiderList(c *gin.Context) {
	var page forms.PageForm

	if err := c.ShouldBindQuery(&page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if total, spiders, err := services.QuerySpiderPage(page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccessList(c, total, spiders)
}

func GetSpider(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	if spider, err := services.QuerySpider(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if spider == nil {
		HandleError(http.StatusNotFound, c, errors.New("spider not found"))
		return
	}

	HandleSuccess(c, spider)
}

func CreateSpider(c *gin.Context) {
	var form forms.SpiderForm

	if err := c.ShouldBind(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 正则校验爬虫名称
	if ok, err := regexp.MatchString("[\\w_-]", form.Name); err != nil || ok == false {
		HandleError(http.StatusBadRequest, c, errors.New("invalid spider name"))
		return
	}

	// 如果不为 zip 文件，返回错误
	if !strings.HasSuffix(form.File.Filename, ".zip") {
		HandleError(http.StatusBadRequest, c, errors.New("invalid zip file"))
		return
	}

	if res, err := services.AddSpider(form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if _, err2 := services.AddSpiderVersion(res.Id, form.SpiderUploadForm); err2 != nil {
		HandleError(http.StatusBadRequest, c, err2)
		return
	}

	HandleSuccess(c, res)
}

func DeleteSpider(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	if res, err := services.RemoveSpider(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}

func GetSpiderVersionList(c *gin.Context) {
	var page forms.SpiderVersionPageForm

	if err := c.ShouldBindUri(&page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if _, err := uuid.FromString(page.SpiderId); err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid spider id"))
		return
	}

	if total, versions, err := services.QuerySpiderVersionPage(page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccessList(c, total, versions)
}

func CreateSpiderVersion(c *gin.Context) {
	var form forms.SpiderUploadForm

	if err := c.ShouldBind(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 如果不为 zip 文件，返回错误
	if !strings.HasSuffix(form.File.Filename, ".zip") {
		HandleError(http.StatusBadRequest, c, errors.New("invalid zip file"))
		return
	}

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	if res, err := services.AddSpiderVersion(id, form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}

func DeleteSpiderVersion(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid spider id"))
		return
	}
	versionId, err := uuid.FromString(c.Param("versionId"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid version id"))
		return
	}

	if res, err := services.RemoveSpiderVersion(id, versionId); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}
