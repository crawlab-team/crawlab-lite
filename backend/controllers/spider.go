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

func GetSpiderList(c *gin.Context) {
	var page forms.PageForm

	if err := c.ShouldBindQuery(&page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if total, spiders, err := services.QuerySpiderList(page.PageNum, page.PageSize); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccessList(c, total, spiders)
	}
}

func GetSpider(c *gin.Context) {
	name := c.Param("name")

	if spider, err := services.QuerySpiderByName(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		if spider == nil {
			HandleError(http.StatusNotFound, c, errors.New("spider not found"))
			return
		}
		HandleSuccess(c, spider)
	}
}

func CreateSpider(c *gin.Context) {
	var form forms.SpiderForm

	if err := c.ShouldBindJSON(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 正则校验项目名称
	if ok, err := regexp.MatchString("[\\w_-]", form.Name); err != nil || ok == false {
		HandleError(http.StatusBadRequest, c, errors.New("invalid spider name"))
	}

	if res, err := services.SaveSpider(form); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func DeleteSpider(c *gin.Context) {
	name := c.Param("name")

	if res, err := services.RemoveSpider(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func GetSpiderVersionList(c *gin.Context) {
	name := c.Param("name")

	if res, err := services.QuerySpiderVersionList(name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func UploadSpiderVersion(c *gin.Context) {
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

	spiderName := c.Param("name")
	if res, err := services.SaveSpiderVersion(spiderName, form); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func DeleteSpiderVersion(c *gin.Context) {
	name := c.Param("name")
	versionId := c.Param("versionId")

	if res, err := services.RemoveSpiderVersion(name, versionId); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}
