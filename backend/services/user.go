package services

import (
	"crawlab-lite/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

func CheckUser(username string, password string) (bool, error) {
	user, err := model.GetUserByName(username)
	if err != nil {
		return false, err
	}
	return user != nil && password == user.Password, nil
}

func MakeToken(username string) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"nbf":      time.Now().Unix(),
	})
	return token.SignedString([]byte(viper.GetString("server.secret")))
}

func GetUserFromToken(tokenStr string) (username string, err error) {
	token, err := jwt.Parse(tokenStr, getSecret())
	if err != nil {
		return "", err
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return "", err
	}

	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return "", err
	}

	username = claim["username"].(string)
	if model.ExistUser(username) == false {
		err = errors.New("username does not match")
		return "", err
	}

	return username, nil
}

func getSecret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("server.secret")), nil
	}
}
