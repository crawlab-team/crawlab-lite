package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type SettingBody struct {
	AllowRegister     string `json:"allow_register"`
	EnableTutorial    string `json:"enable_tutorial"`
	RunOnMaster       string `json:"run_on_master"`
	EnableDemoSpiders string `json:"enable_demo_spiders"`
}

func GetVersion(c *gin.Context) {
	version := viper.GetString("version")

	HandleSuccess(c, version)
}

func GetSetting(c *gin.Context) {
	body := SettingBody{
		EnableTutorial:    viper.GetString("setting.enableTutorial"),
		RunOnMaster:       viper.GetString("setting.runOnMaster"),
		EnableDemoSpiders: viper.GetString("setting.enableDemoSpiders"),
	}

	HandleSuccess(c, body)
}
