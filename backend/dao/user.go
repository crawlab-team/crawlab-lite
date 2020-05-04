package dao

import (
	"crawlab-lite/models"
	"github.com/spf13/viper"
)

func GetUserByName(username string) (user *models.User, err error) {
	userList := make([]*models.User, 0)
	if err := viper.UnmarshalKey("users", &userList); err != nil {
		return nil, err
	}
	for _, user := range userList {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func GetUserList() ([]*models.User, error) {
	userList := make([]*models.User, 0)
	if err := viper.UnmarshalKey("users", &userList); err != nil {
		return nil, err
	}
	return userList, nil
}

func ExistUser(username string) bool {
	user, err := GetUserByName(username)
	return user != nil && err == nil
}
