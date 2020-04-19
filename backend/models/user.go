package models

import "github.com/spf13/viper"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserByName(username string) (*User, error) {
	userList := make([]*User, 0)
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

func GetUserList() ([]*User, error) {
	userList := make([]*User, 0)
	if err := viper.UnmarshalKey("users", &userList); err != nil {
		return nil, err
	}
	return userList, nil
}

func ExistUser(username string) bool {
	user, err := GetUserByName(username)
	return user != nil && err == nil
}
