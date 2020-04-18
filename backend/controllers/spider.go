package controllers

import (
	"crawlab-lite/forms"
	"crawlab-lite/model"
	"crawlab-lite/services"
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func QuerySpiderList(c *gin.Context) {
	spiders, err := services.GetSpiderList()
	if err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	HandleSuccess(c, spiders)
}

func QuerySpider(c *gin.Context) {
	name := c.Param("name")

	spider, err := services.GeySpiderByName(name)
	if err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if spider == nil {
		HandleError(http.StatusNotFound, c, errors.New("spider not found"))
		return
	}

	HandleSuccess(c, spider)
}

func UploadSpider(c *gin.Context) {
	var form forms.SpiderUploadForm

	if err := c.ShouldBind(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 正则校验爬虫名称
	if ok, err := regexp.MatchString("[\\w_-]", form.Name); err != nil || ok == false {
		HandleError(http.StatusBadRequest, c, errors.New("invalid spider name"))
	}

	// 如果不为 zip 文件，返回错误
	if !strings.HasSuffix(form.File.Filename, ".zip") {
		HandleError(http.StatusBadRequest, c, errors.New("invalid zip file"))
		return
	}

	// 生成爬虫文件名
	fileName := uuid.NewV4().String() + ".zip"

	spider := &model.Spider{
		Name: form.Name,
		Cmd:  form.Cmd,
		File: fileName,
	}

	// 保存爬虫文件
	if err := os.MkdirAll(spider.GetSpiderDirPath(), os.ModePerm); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	if err := c.SaveUploadedFile(form.File, spider.GetSpiderFilePath()); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 保存爬虫信息
	if err := spider.SaveJson(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
}

func DeleteSpider(c *gin.Context) {
	name := c.Param("name")

	data, err := services.DeleteSpiderByName(name)
	if err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, data)
}
