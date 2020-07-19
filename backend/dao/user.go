package dao

import (
	"crawlab-lite/models"
	"github.com/spf13/viper"
)

func GetUser() (user *models.User) {
	username := viper.GetString("user.username")
	password := viper.GetString("user.password")
	if username == "" || password == "" {
		return nil
	}
	return &models.User{
		Username: username,
		Password: password,
	}
}

func ExistUser(username string) bool {
	return username != "" && username == viper.GetString("user.username")
}
