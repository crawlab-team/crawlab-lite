package controllers

import (
	"crawlab-lite/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func GetLatestRelease(c *gin.Context) {
	latestRelease, err := services.GetLatestRelease()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
	}
	HandleSuccess(c, latestRelease)
}
